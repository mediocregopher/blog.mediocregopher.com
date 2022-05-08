package post

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"sync"
)

var (
	// ErrAssetNotFound is used to indicate an Asset could not be found in the
	// AssetStore.
	ErrAssetNotFound = errors.New("asset not found")
)

// AssetStore implements the storage and retrieval of binary assets, which are
// intended to be used by posts (e.g. images).
type AssetStore interface {

	// Set sets the id to the contents of the given io.Reader.
	Set(id string, from io.Reader) error

	// Get writes the id's body to the given io.Writer, or returns
	// ErrAssetNotFound.
	Get(id string, into io.Writer) error

	// Delete's the body stored for the id, if any.
	Delete(id string) error
}

type assetStore struct {
	db *sql.DB
}

// NewAssetStore initializes a new AssetStore using an existing SQLDB.
func NewAssetStore(db *SQLDB) AssetStore {
	return &assetStore{
		db: db.db,
	}
}

func (s *assetStore) Set(id string, from io.Reader) error {

	body, err := io.ReadAll(from)
	if err != nil {
		return fmt.Errorf("reading body fully into memory: %w", err)
	}

	_, err = s.db.Exec(
		`INSERT INTO assets (id, body)
		VALUES (?, ?)
		ON CONFLICT (id) DO UPDATE SET body=excluded.body`,
		id, body,
	)

	if err != nil {
		return fmt.Errorf("inserting into assets: %w", err)
	}

	return nil
}

func (s *assetStore) Get(id string, into io.Writer) error {

	var body []byte

	err := s.db.QueryRow(`SELECT body FROM assets WHERE id = ?`, id).Scan(&body)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrAssetNotFound
	} else if err != nil {
		return fmt.Errorf("selecting from assets: %w", err)
	}

	if _, err := io.Copy(into, bytes.NewReader(body)); err != nil {
		return fmt.Errorf("writing body to io.Writer: %w", err)
	}

	return nil
}

func (s *assetStore) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM assets WHERE id = ?`, id)
	return err
}

////////////////////////////////////////////////////////////////////////////////

type cachedAssetStore struct {
	inner AssetStore
	m     sync.Map
}

// NewCachedAssetStore wraps an AssetStore in an in-memory cache.
func NewCachedAssetStore(assetStore AssetStore) AssetStore {
	return &cachedAssetStore{
		inner: assetStore,
	}
}

func (s *cachedAssetStore) Set(id string, from io.Reader) error {

	buf := new(bytes.Buffer)
	from = io.TeeReader(from, buf)

	if err := s.inner.Set(id, from); err != nil {
		return err
	}

	s.m.Store(id, buf.Bytes())
	return nil
}

func (s *cachedAssetStore) Get(id string, into io.Writer) error {

	if bodyI, ok := s.m.Load(id); ok {

		if err, ok := bodyI.(error); ok {
			return err
		}

		if _, err := io.Copy(into, bytes.NewReader(bodyI.([]byte))); err != nil {
			return fmt.Errorf("writing body to io.Writer: %w", err)
		}

		return nil
	}

	buf := new(bytes.Buffer)
	into = io.MultiWriter(into, buf)

	if err := s.inner.Get(id, into); errors.Is(err, ErrAssetNotFound) {
		s.m.Store(id, err)
		return err
	} else if err != nil {
		return err
	}

	s.m.Store(id, buf.Bytes())
	return nil
}

func (s *cachedAssetStore) Delete(id string) error {

	if err := s.inner.Delete(id); err != nil {
		return err
	}

	s.m.Delete(id)
	return nil
}

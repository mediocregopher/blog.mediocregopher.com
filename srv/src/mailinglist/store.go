package mailinglist

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

var (
	// ErrNotFound is used to indicate an email could not be found in the
	// database.
	ErrNotFound = errors.New("no record found")
)

// EmailIterator will iterate through a sequence of emails, returning the next
// email in the sequence on each call, or returning io.EOF.
type EmailIterator func() (Email, error)

// Email describes all information related to an email which has yet
// to be verified.
type Email struct {
	Email     string
	SubToken  string
	CreatedAt time.Time

	UnsubToken string
	VerifiedAt time.Time
}

// Store is used for storing MailingList related information.
type Store interface {

	// Set is used to set the information related to an email.
	Set(Email) error

	// Get will return the record for the given email, or ErrNotFound.
	Get(email string) (Email, error)

	// GetBySubToken will return the record for the given SubToken, or
	// ErrNotFound.
	GetBySubToken(subToken string) (Email, error)

	// GetByUnsubToken will return the record for the given UnsubToken, or
	// ErrNotFound.
	GetByUnsubToken(unsubToken string) (Email, error)

	// Delete will delete the record for the given email.
	Delete(email string) error

	// GetAll returns all emails for which there is a record.
	GetAll() EmailIterator

	Close() error
}

var migrations = []*migrate.Migration{
	&migrate.Migration{
		Id: "1",
		Up: []string{
			`CREATE TABLE emails (
				id          TEXT PRIMARY KEY,
				email       TEXT NOT NULL,
				sub_token   TEXT NOT NULL,
				created_at  INTEGER NOT NULL,

				unsub_token TEXT,
				verified_at INTEGER
			)`,
		},
		Down: []string{"DROP TABLE emails"},
	},
}

type store struct {
	db *sql.DB
}

// NewStore initializes a new Store using a sqlite3 database at the given file
// path.
func NewStore(dbFile string) (Store, error) {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite file: %w", err)
	}

	migrations := &migrate.MemoryMigrationSource{Migrations: migrations}

	if _, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return &store{
		db: db,
	}, nil
}

func (s *store) emailID(email string) string {
	email = strings.ToLower(email)
	h := sha512.New()
	h.Write([]byte(email))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (s *store) Set(email Email) error {
	_, err := s.db.Exec(
		`INSERT INTO emails (
			id, email, sub_token, created_at, unsub_token, verified_at
		)
		VALUES
		(?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			email=excluded.email,
			sub_token=excluded.sub_token,
			unsub_token=excluded.unsub_token,
			verified_at=excluded.verified_at
		`,
		s.emailID(email.Email),
		email.Email,
		email.SubToken,
		email.CreatedAt.Unix(),
		email.UnsubToken,
		sql.NullInt64{
			Int64: email.VerifiedAt.Unix(),
			Valid: !email.VerifiedAt.IsZero(),
		},
	)

	return err
}

var scanCols = []string{
	"email", "sub_token", "created_at", "unsub_token", "verified_at",
}

type row interface {
	Scan(...interface{}) error
}

func (s *store) scanRow(row row) (Email, error) {
	var email Email
	var createdAt int64
	var verifiedAt sql.NullInt64

	err := row.Scan(
		&email.Email,
		&email.SubToken,
		&createdAt,
		&email.UnsubToken,
		&verifiedAt,
	)
	if err != nil {
		return Email{}, err
	}

	email.CreatedAt = time.Unix(createdAt, 0)
	if verifiedAt.Valid {
		email.VerifiedAt = time.Unix(verifiedAt.Int64, 0)
	}

	return email, nil
}

func (s *store) scanSingleRow(row *sql.Row) (Email, error) {
	email, err := s.scanRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Email{}, ErrNotFound
	}

	return email, err
}

func (s *store) Get(email string) (Email, error) {
	row := s.db.QueryRow(
		`SELECT `+strings.Join(scanCols, ",")+`
		FROM emails
		WHERE id=?`,
		s.emailID(email),
	)

	return s.scanSingleRow(row)
}

func (s *store) GetBySubToken(subToken string) (Email, error) {
	row := s.db.QueryRow(
		`SELECT `+strings.Join(scanCols, ",")+`
		FROM emails
		WHERE sub_token=?`,
		subToken,
	)

	return s.scanSingleRow(row)
}

func (s *store) GetByUnsubToken(unsubToken string) (Email, error) {
	row := s.db.QueryRow(
		`SELECT `+strings.Join(scanCols, ",")+`
		FROM emails
		WHERE unsub_token=?`,
		unsubToken,
	)

	return s.scanSingleRow(row)
}

func (s *store) Delete(email string) error {
	_, err := s.db.Exec(
		`DELETE FROM emails WHERE id=?`,
		s.emailID(email),
	)
	return err
}

func (s *store) GetAll() EmailIterator {
	rows, err := s.db.Query(
		`SELECT ` + strings.Join(scanCols, ",") + `
		FROM emails`,
	)

	return func() (Email, error) {
		if err != nil {
			return Email{}, err

		} else if !rows.Next() {
			return Email{}, io.EOF
		}
		return s.scanRow(rows)
	}
}

func (s *store) Close() error {
	return s.db.Close()
}

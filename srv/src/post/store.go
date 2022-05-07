package post

import (
	"database/sql"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // we need dis
	migrate "github.com/rubenv/sql-migrate"
)

var (
	// ErrNotFound is used to indicate a Post could not be found in the
	// database.
	ErrNotFound = errors.New("not found")
)

// StoredPost is a Post which has been stored in a Store, and has been given
// some extra fields as a result.
type StoredPost struct {
	Post

	PublishedAt   time.Time
	LastUpdatedAt time.Time
}

// URL returns the relative URL of the StoredPost.
func (p StoredPost) URL() string {
	return path.Join(
		fmt.Sprintf(
			"%d/%0d/%0d",
			p.PublishedAt.Year(),
			p.PublishedAt.Month(),
			p.PublishedAt.Day(),
		),
		p.ID+".html",
	)
}

// Store is used for storing posts to a persistent storage.
type Store interface {

	// Set sets the Post data into the storage, keyed by the Post's ID. It
	// overwrites a previous Post with the same ID, if there was one.
	Set(post Post, now time.Time) error

	// Get returns count StoredPosts, sorted time descending, offset by the given page
	// number. The returned boolean indicates if there are more pages or not.
	Get(page, count int) ([]StoredPost, bool, error)

	// GetByID will return the StoredPost with the given ID, or ErrNotFound.
	GetByID(id string) (StoredPost, error)

	// GetBySeries returns all StoredPosts with the given series, sorted time
	// ascending, or empty slice.
	GetBySeries(series string) ([]StoredPost, error)

	// GetByTag returns all StoredPosts with the given tag, sorted time
	// ascending, or empty slice.
	GetByTag(tag string) ([]StoredPost, error)

	// Delete will delete the StoredPost with the given ID.
	Delete(id string) error

	Close() error
}

var migrations = []*migrate.Migration{
	&migrate.Migration{
		Id: "1",
		Up: []string{
			`CREATE TABLE posts (
				id          TEXT NOT NULL PRIMARY KEY,
				title       TEXT NOT NULL,
				description TEXT NOT NULL,
				series      TEXT,

				published_at    INTEGER NOT NULL,
				last_updated_at INTEGER,

				body TEXT NOT NULL
			)`,
			`CREATE TABLE post_tags (
				post_id TEXT NOT NULL,
				tag     TEXT NOT NULL,
				UNIQUE(post_id, tag)
			)`,
		},
		Down: []string{
			"DROP TABLE post_tags",
			"DROP TABLE posts",
		},
	},
}

// Params are parameters used to initialize a new Store. All fields are required
// unless otherwise noted.
type StoreParams struct {

	// Path to the file the database will be stored at.
	DBFilePath string
}

type store struct {
	params StoreParams
	db     *sql.DB
}

// NewStore initializes a new Store using a sqlite3 database at the given file
// path.
func NewStore(params StoreParams) (Store, error) {

	db, err := sql.Open("sqlite3", params.DBFilePath)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite file: %w", err)
	}

	migrations := &migrate.MemoryMigrationSource{Migrations: migrations}

	if _, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return &store{
		params: params,
		db:     db,
	}, nil
}

func (s *store) Close() error {
	return s.db.Close()
}

// if the callback returns an error then the transaction is aborted.
func (s *store) withTx(cb func(*sql.Tx) error) error {

	tx, err := s.db.Begin()

	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if err := cb(tx); err != nil {

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf(
				"rolling back transaction: %w (original error: %v)",
				rollbackErr, err,
			)
		}

		return fmt.Errorf("performing transaction: %w (rolled back)", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (s *store) Set(post Post, now time.Time) error {
	return s.withTx(func(tx *sql.Tx) error {

		nowTS := now.Unix()

		nowSql := sql.NullInt64{Int64: nowTS, Valid: !now.IsZero()}

		_, err := tx.Exec(
			`INSERT INTO posts (
			id, title, description, series, published_at, body
		)
		VALUES
		(?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			title=excluded.title,
			description=excluded.description,
			series=excluded.series,
			last_updated_at=?,
			body=excluded.body`,
			post.ID,
			post.Title,
			post.Description,
			&sql.NullString{String: post.Series, Valid: post.Series != ""},
			nowSql,
			post.Body,
			nowSql,
		)

		if err != nil {
			return fmt.Errorf("inserting into posts: %w", err)
		}

		// this is a bit of a hack, but it allows us to update the tagset without
		// doing a diff.
		_, err = tx.Exec(`DELETE FROM post_tags WHERE post_id = ?`, post.ID)

		if err != nil {
			return fmt.Errorf("clearning post tags: %w", err)
		}

		for _, tag := range post.Tags {

			_, err = tx.Exec(
				`INSERT INTO post_tags (post_id, tag) VALUES (?, ?)
			ON CONFLICT DO NOTHING`,
				post.ID,
				tag,
			)

			if err != nil {
				return fmt.Errorf("inserting tag %q: %w", tag, err)
			}
		}

		return nil
	})
}

func (s *store) get(
	querier interface {
		Query(string, ...interface{}) (*sql.Rows, error)
	},
	limit, offset int,
	where string, whereArgs ...interface{},
) (
	[]StoredPost, error,
) {

	query := `SELECT
		p.id, p.title, p.description, p.series, pt.tag,
		p.published_at, p.last_updated_at,
		p.body
	FROM posts p
	LEFT JOIN post_tags pt ON (p.id = pt.post_id)
	` + where + `
	ORDER BY p.published_at ASC, p.title ASC`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	rows, err := querier.Query(query, whereArgs...)

	if err != nil {
		return nil, fmt.Errorf("selecting: %w", err)
	}

	var posts []StoredPost

	for rows.Next() {

		var (
			post                       StoredPost
			series, tag                sql.NullString
			publishedAt, lastUpdatedAt sql.NullInt64
		)

		err := rows.Scan(
			&post.ID, &post.Title, &post.Description, &series, &tag,
			&publishedAt, &lastUpdatedAt,
			&post.Body,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		if tag.Valid {

			if l := len(posts); l > 0 && posts[l-1].ID == post.ID {
				posts[l-1].Tags = append(posts[l-1].Tags, tag.String)
				continue
			}

			post.Tags = append(post.Tags, tag.String)
		}

		post.Series = series.String

		if publishedAt.Valid {
			post.PublishedAt = time.Unix(publishedAt.Int64, 0).UTC()
		}

		if lastUpdatedAt.Valid {
			post.LastUpdatedAt = time.Unix(lastUpdatedAt.Int64, 0).UTC()
		}

		posts = append(posts, post)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("closing row iterator: %w", err)
	}

	return posts, nil
}

func (s *store) Get(page, count int) ([]StoredPost, bool, error) {

	posts, err := s.get(s.db, count+1, page*count, ``)

	if err != nil {
		return nil, false, fmt.Errorf("querying posts: %w", err)
	}

	var hasMore bool

	if len(posts) > count {
		hasMore = true
		posts = posts[:count]
	}

	return posts, hasMore, nil
}

func (s *store) GetByID(id string) (StoredPost, error) {

	posts, err := s.get(s.db, 0, 0, `WHERE p.id=?`, id)

	if err != nil {
		return StoredPost{}, fmt.Errorf("querying posts: %w", err)
	}

	if len(posts) == 0 {
		return StoredPost{}, ErrNotFound
	}

	if len(posts) > 1 {
		panic(fmt.Sprintf("got back multiple posts querying id %q: %+v", id, posts))
	}

	return posts[0], nil
}

func (s *store) GetBySeries(series string) ([]StoredPost, error) {
	return s.get(s.db, 0, 0, `WHERE p.series=?`, series)
}

func (s *store) GetByTag(tag string) ([]StoredPost, error) {

	var posts []StoredPost

	err := s.withTx(func(tx *sql.Tx) error {

		rows, err := tx.Query(`SELECT post_id FROM post_tags WHERE tag = ?`, tag)

		if err != nil {
			return fmt.Errorf("querying post_tags by tag: %w", err)
		}

		var (
			placeholders []string
			whereArgs    []interface{}
		)

		for rows.Next() {

			var id string

			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return fmt.Errorf("scanning id: %w", err)
			}

			whereArgs = append(whereArgs, id)
			placeholders = append(placeholders, "?")
		}

		if err := rows.Close(); err != nil {
			return fmt.Errorf("closing row iterator: %w", err)
		}

		where := fmt.Sprintf("WHERE p.id IN (%s)", strings.Join(placeholders, ","))

		if posts, err = s.get(tx, 0, 0, where, whereArgs...); err != nil {
			return fmt.Errorf("querying for ids %+v: %w", whereArgs, err)
		}

		return nil
	})

	return posts, err
}

func (s *store) Delete(id string) error {

	tx, err := s.db.Begin()

	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM post_tags WHERE post_id = ?`, id); err != nil {
		return fmt.Errorf("deleting from post_tags: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM posts WHERE id = ?`, id); err != nil {
		return fmt.Errorf("deleting from posts: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

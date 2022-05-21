// Package post deals with the storage and rendering of blog posts.
package post

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	// ErrPostNotFound is used to indicate a Post could not be found in the
	// Store.
	ErrPostNotFound = errors.New("post not found")
)

var titleCleanRegexp = regexp.MustCompile(`[^a-z ]`)

// NewID generates a (hopefully) unique ID based on the given title.
func NewID(title string) string {
	title = strings.ToLower(title)
	title = titleCleanRegexp.ReplaceAllString(title, "")
	title = strings.ReplaceAll(title, " ", "-")
	return title
}

// Post contains all information having to do with a blog post.
type Post struct {
	ID          string
	Title       string
	Description string
	Tags        []string // only alphanumeric supported
	Series      string
	Body        string
}

// StoredPost is a Post which has been stored in a Store, and has been given
// some extra fields as a result.
type StoredPost struct {
	Post

	PublishedAt   time.Time
	LastUpdatedAt time.Time
}

// Store is used for storing posts to a persistent storage.
type Store interface {

	// Set sets the Post data into the storage, keyed by the Post's ID. If there
	// was not a previously existing Post with the same ID then Set returns
	// true. It overwrites the previous Post with the same ID otherwise.
	Set(post Post, now time.Time) (bool, error)

	// Get returns count StoredPosts, sorted time descending, offset by the
	// given page number. The returned boolean indicates if there are more pages
	// or not.
	Get(page, count int) ([]StoredPost, bool, error)

	// GetByID will return the StoredPost with the given ID, or ErrPostNotFound.
	GetByID(id string) (StoredPost, error)

	// GetBySeries returns all StoredPosts with the given series, sorted time
	// descending, or empty slice.
	GetBySeries(series string) ([]StoredPost, error)

	// GetByTag returns all StoredPosts with the given tag, sorted time
	// descending, or empty slice.
	GetByTag(tag string) ([]StoredPost, error)

	// GetTags returns all tags which have at least one Post using them.
	GetTags() ([]string, error)

	// Delete will delete the StoredPost with the given ID.
	Delete(id string) error
}

type store struct {
	db *sql.DB
}

// NewStore initializes a new Store using an existing SQLDB.
func NewStore(db *SQLDB) Store {
	return &store{
		db: db.db,
	}
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

func (s *store) Set(post Post, now time.Time) (bool, error) {

	if post.ID == "" {
		return false, errors.New("post ID can't be empty")
	}

	var first bool

	err := s.withTx(func(tx *sql.Tx) error {

		nowTS := now.Unix()

		nowSQL := sql.NullInt64{Int64: nowTS, Valid: !now.IsZero()}

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
			nowSQL,
			post.Body,
			nowSQL,
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

		err = tx.QueryRow(
			`SELECT 1 FROM posts WHERE id=? AND last_updated_at IS NULL`,
			post.ID,
		).Scan(new(int))

		first = !errors.Is(err, sql.ErrNoRows)

		return nil
	})

	return first, err
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

	query := `
		SELECT
			p.id, p.title, p.description, p.series, GROUP_CONCAT(pt.tag),
			p.published_at, p.last_updated_at,
			p.body
		FROM posts p
		LEFT JOIN post_tags pt ON (p.id = pt.post_id)
		` + where + `
		GROUP BY (p.id)
		ORDER BY p.published_at DESC, p.title DESC`

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

		post.Series = series.String

		if tag.String != "" {
			post.Tags = strings.Split(tag.String, ",")
		}

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
		return StoredPost{}, ErrPostNotFound
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

func (s *store) GetTags() ([]string, error) {

	rows, err := s.db.Query(`SELECT tag FROM post_tags GROUP BY tag`)
	if err != nil {
		return nil, fmt.Errorf("querying all tags: %w", err)
	}
	defer rows.Close()

	var tags []string

	for rows.Next() {

		var tag string

		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("scanning tag: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
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

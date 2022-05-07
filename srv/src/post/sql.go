package post

import (
	"database/sql"
	"fmt"
	"path"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	migrate "github.com/rubenv/sql-migrate"

	_ "github.com/mattn/go-sqlite3" // we need dis
)

var migrations = &migrate.MemoryMigrationSource{Migrations: []*migrate.Migration{
	{
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

			`CREATE TABLE assets (
				id   TEXT NOT NULL PRIMARY KEY,
				body BLOB NOT NULL
			)`,
		},
		Down: []string{
			"DROP TABLE assets",
			"DROP TABLE post_tags",
			"DROP TABLE posts",
		},
	},
}}

// SQLDB is a sqlite3 database which can be used by storage interfaces within
// this package.
type SQLDB struct {
	db *sql.DB
}

// NewSQLDB initializes and returns a new sqlite3 database for storage
// intefaces. The db will  be created within the given data directory.
func NewSQLDB(dataDir cfg.DataDir) (*SQLDB, error) {

	path := path.Join(dataDir.Path, "post.sqlite3")

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite file at %q: %w", path, err)
	}

	if _, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return &SQLDB{db}, nil
}

// NewSQLDB is like NewSQLDB, but the database will be initialized in memory.
func NewInMemSQLDB() *SQLDB {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(fmt.Errorf("opening sqlite in memory: %w", err))
	}

	if _, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up); err != nil {
		panic(fmt.Errorf("running migrations: %w", err))
	}

	return &SQLDB{db}
}

// Close cleans up loose resources being held by the db.
func (db *SQLDB) Close() error {
	return db.db.Close()
}

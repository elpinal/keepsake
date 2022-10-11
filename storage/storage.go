package storage

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/elpinal/keepsake/entry"
	"github.com/elpinal/keepsake/log"
)

type SQLite3Storage struct {
	logger log.Logger
	*sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS entries (id integer NOT NULL PRIMARY KEY, url text NOT NULL, title text, date timestamp NOT NULL);
`

func New(logger log.Logger, path string) (*SQLite3Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &SQLite3Storage{
		logger: logger,
		DB:     db,
	}, nil
}

func (db *SQLite3Storage) Add(url string, title string, date time.Time) error {
	_, err := db.Exec(
		`INSERT INTO entries (url, title, date) VALUES (?, ?, ?)`,
		url,
		title,
		date,
	)
	if err != nil {
		return err
	}
	db.logger.LogInfo("sqlite3: inserted an item", url)
	return nil
}

func (db *SQLite3Storage) List(limit int, offset int) ([]entry.Entry, error) {
	db.logger.LogInfo("sqlite3: select", map[string]int{"limit": limit, "offset": offset})
	rows, err := db.Query(
		`SELECT url, title, date FROM entries ORDER BY id DESC LIMIT ?, ?`,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]entry.Entry, 0, limit)

	for rows.Next() {
		var (
			url   string
			title string
			date  time.Time
		)
		if err := rows.Scan(&url, &title, &date); err != nil {
			return nil, err
		}
		entries = append(entries, entry.Entry{URL: url, Title: title, Date: date})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (db *SQLite3Storage) Count() (int, error) {
	row := db.QueryRow(`SELECT COUNT(*) FROM entries`)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *SQLite3Storage) Export(enc *json.Encoder, limit int, offset int) error {
	db.logger.LogInfo("sqlite3: select", map[string]int{"limit": limit, "offset": offset})
	rows, err := db.Query(
		`SELECT url, title, date FROM entries ORDER BY id ASC LIMIT ?, ?`,
		offset,
		limit,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			url   string
			title string
			date  time.Time
		)
		if err := rows.Scan(&url, &title, &date); err != nil {
			return err
		}
		err = enc.Encode(entry.Entry{URL: url, Title: title, Date: date})
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (db *SQLite3Storage) Import(dec *json.Decoder) error {
	// TODO: timeout or something?

	stmt, err := db.Prepare(`INSERT INTO entries (url, title, date) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for dec.More() {
		var entry entry.Entry
		if err := dec.Decode(&entry); err != nil {
			return err
		}
		db.logger.LogInfo("importing an entry", entry)
		_, err = stmt.Exec(entry.URL, entry.Title, entry.Date)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *SQLite3Storage) Search(q string) ([]entry.Entry, error) {
	db.logger.LogInfo("sqlite3: query", q)
	rows, err := db.Query(`SELECT url, title, date FROM entries WHERE title LIKE ? ORDER BY id DESC`, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]entry.Entry, 0)

	for rows.Next() {
		var (
			url   string
			title string
			date  time.Time
		)
		if err := rows.Scan(&url, &title, &date); err != nil {
			return nil, err
		}
		entries = append(entries, entry.Entry{URL: url, Title: title, Date: date})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

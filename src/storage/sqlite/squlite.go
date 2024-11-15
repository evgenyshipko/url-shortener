package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log/slog"
	"url-shortener/src/storage"
)

type Storage struct {
	db *sql.DB
}

// TODO: перейти на postgres

func New(storagePath string) (*Storage, error) {
	// FIXME: почему не рефлексией?
	const funcName = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	// TODO: добавить миграции
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(url string, alias string) (int64, error) {
	const funcName = "storage.sqlite.SaveURL"

	statement, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", funcName, err)
	}

	res, err := statement.Exec(url, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", funcName, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", funcName, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", funcName, storage.ErrURLExists)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const funcName = "storage.sqlite.GetURL"

	statement, err := s.db.Prepare("SELECT url FROM url WHERE alias = (?)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", funcName, err)
	}

	var resURL string
	err = statement.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", funcName, err)
	}

	slog.Info("ded!", resURL)

	return resURL, nil
}

// TODO: implement method
// func (s *Storage) DeleteUrl(alias string) error

package sqlite

import (
  "database/sql"
  "errors"
  "fmt"
  "github.com/mattn/go-sqlite3"
  "go-url-shortner/internal/storage"
)

type Storage struct {
  db *sql.DB
}

func New(storagePath string) (*Storage, error) {
  const op = "storage.sqlite.New" // там мы враппим ошибку и опказываем где она произошла

  db, err := sql.Open("sqlite3", storagePath)
  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS url(
        id INTEGER PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL,
        created_at TEXT NOT NULL DEFAULT (datetime('now')),
    		updated_at TEXT NOT NULL DEFAULT (datetime('now'))
    );
    CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
  `)
  if err != nil {
    return nil, fmt.Errorf("#{op}: #{err}")
  }

  _, err = stmt.Exec()
  if err != nil {
    return nil, fmt.Errorf("%s: %w", op, err)
  }

  return &Storage{db: db}, nil
}

func (s *Storage) Save(urlToSave string, alias string) (int64, error) {
  const op = "storage.sqlite.SaveURL"

  smtt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
  if err != nil {
    return 0, fmt.Errorf("%s: %w", op, err)
  }

  res, err := smtt.Exec(urlToSave, alias)
  if err != nil {
    if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
      return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
    }
    return 0, fmt.Errorf("%s: %w", op, err)
  }

  id, err := res.LastInsertId()
  if err != nil {
    return 0, fmt.Errorf("%s: failde to insert id: %w", op, err)
  }

  return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
  const op = "storage.sqlite.GetURL"

  stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
  if err != nil {
    return "", fmt.Errorf("%s: %w", op, err)
  }

  var resURL string
  err = stmt.QueryRow(alias).Scan(&resURL)
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
    }
    return "", fmt.Errorf("%s: %w", op, err)
  }

  return resURL, nil
}

// todo check DELTE
func (s *Storage) DeleteURL(alias string) error {
  const op = "storage.sqlite.GetURL"

  stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
  if err != nil {
    return fmt.Errorf("%s: %w", op, err)
  }

  err = stmt.QueryRow(alias).Err()
  return err
}

package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./urls.db")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hash TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL UNIQUE
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetURLByHash(db *sql.DB, hash string) (string, error) {
	var result string
	err := db.QueryRow("SELECT url FROM urls WHERE hash = ?", hash).Scan(&result)
	if err != nil {
		return "", nil
	}

	return result, nil
}

func CreateURL(db *sql.DB, hash string, url string) error {
	result, err := db.Exec("INSERT INTO urls (hash, url) VALUES (?, ?)", hash, url)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

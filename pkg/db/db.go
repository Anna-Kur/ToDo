package db

import (
	"database/sql"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	dbPath := os.Getenv("TODO_DBFILE")

	if dbPath == "" {
		dbPath = dbFile
	}

	_, err := os.Stat(dbPath)
	install := os.IsNotExist(err)

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return errors.New("Error opening database: " + err.Error())
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return errors.New("Table and index creation error: " + err.Error())
		}
	}

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

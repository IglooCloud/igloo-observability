package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Exec(db *sql.DB, query string, args ...any) (sql.Result, error) {
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil, err
	}

	return stmt.Exec(args...)
}

func Query(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil, err
	}

	return stmt.Query(args...)
}

func Initialize(db *sql.DB) error {
	// create table
	sql_table := `
	CREATE TABLE IF NOT EXISTS realrecords(
		value REAL NOT NULL,
		timestamp INTEGER NOT NULL,
		bucket TEXT NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_realrecords_bucket_timestamp ON realrecords(bucket, timestamp);
	`

	_, err := Exec(db, sql_table)
	return err
}

func RealFromDB(db *sql.DB) *Real {
	return &Real{db}
}

func Connect(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = Initialize(db)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return db, nil
}

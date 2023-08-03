package storage

import (
	"database/sql"
	"fmt"

	"github.com/IglooCloud/igloo-observability/internal/log"

	_ "github.com/mattn/go-sqlite3"
)

var logger = log.Default().Service("store")

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
	_, err := Exec(db, `CREATE TABLE IF NOT EXISTS realrecords(
		value REAL NOT NULL,
		timestamp INTEGER NOT NULL,
		bucket TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}

	_, err = Exec(db, `CREATE INDEX IF NOT EXISTS idx_realrecords_bucket_timestamp ON realrecords(bucket, timestamp);`)
	if err != nil {
		return err
	}

	_, err = Exec(db, `CREATE TABLE IF NOT EXISTS intperiodrecords(
		value REAL NOT NULL,
		start INTEGER NOT NULL,
		end INTEGER NOT NULL,
		bucket TEXT NOT NULL
	);`)
	if err != nil {
		return err
	}

	_, err = Exec(db, `CREATE INDEX IF NOT EXISTS idx_intperiodrecords_bucket_timestamp ON intperiodrecords(bucket, start, end);`)
	if err != nil {
		return err
	}

	return err
}

func RealFromDB(db *sql.DB) *Real {
	return &Real{db}
}

func PeriodIntFromDB(db *sql.DB) *PeriodInt {
	return &PeriodInt{db}
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

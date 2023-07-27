package storage

import (
	"database/sql"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/warehouse"
)

type Real struct {
	db *sql.DB
}

func (r *Real) Write(bucket string, value warehouse.Record[float64]) error {
	_, err := Exec(
		r.db,
		"INSERT INTO realrecords(value, timestamp, bucket) values(?, ?, ?)",
		value.Value,
		value.Timestamp.Unix(),
		bucket,
	)
	return err
}

func (r *Real) Read(bucket string, start, end time.Time) (warehouse.RecordIterator[float64], error) {
	rows, err := Query(
		r.db,
		"SELECT value, timestamp FROM realrecords WHERE timestamp >= ? AND timestamp <= ? AND bucket = ?",
		start.Unix(),
		end.Unix(),
		bucket,
	)
	if err != nil {
		return nil, err
	}

	return RealRecordIterator{rows}, nil
}

type RealRecordIterator struct {
	rows *sql.Rows
}

func (r RealRecordIterator) Next() (warehouse.Record[float64], bool) {
	r.rows.Next()

	var value warehouse.Record[float64]
	var rawTimestamp int64
	err := r.rows.Scan(&value.Value, &rawTimestamp)
	if err != nil {
		return value, false
	}

	value.Timestamp = time.Unix(rawTimestamp, 0)

	return value, true
}

func (r RealRecordIterator) Err() error {
	return r.rows.Err()
}

func (r RealRecordIterator) Close() error {
	return r.rows.Close()
}

func (r *Real) ReadValues(bucket string, start, end time.Time) (warehouse.Iterator[float64], error) {
	rows, err := Query(
		r.db,
		"SELECT value FROM realrecords WHERE timestamp >= ? AND timestamp <= ? AND bucket = ?",
		start.Unix(),
		end.Unix(),
		bucket,
	)
	if err != nil {
		return nil, err
	}

	return RealIterator{rows}, nil
}

type RealIterator struct {
	rows *sql.Rows
}

func (r RealIterator) Next() (float64, bool) {
	r.rows.Next()

	var value float64
	err := r.rows.Scan(&value)
	if err != nil {
		return 0, false
	}

	return value, true
}

func (r RealIterator) Err() error {
	return r.rows.Err()
}

func (r RealIterator) Close() error {
	return r.rows.Close()
}

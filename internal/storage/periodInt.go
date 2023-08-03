package storage

import (
	"database/sql"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/warehouse"
)

type PeriodInt struct {
	db *sql.DB
}

func (r *PeriodInt) Write(bucket string, value warehouse.PeriodRecord[int64]) error {
	_, err := Exec(
		r.db,
		"INSERT INTO intperiodrecords(value, start, end, bucket) values(?, ?, ?, ?)",
		value.Value,
		value.Start.Unix(),
		value.End.Unix(),
		bucket,
	)
	return err
}

func (r *PeriodInt) Read(bucket string, start, end time.Time) (warehouse.PeriodRecordIterator[int64], error) {
	rows, err := Query(
		r.db,
		"SELECT value, start, end FROM intperiodrecords WHERE end >= ? AND start <= ? AND bucket = ?",
		start.Unix(),
		end.Unix(),
		bucket,
	)
	if err != nil {
		return nil, err
	}

	return IntPeriodRecordIterator{rows}, nil
}

type IntPeriodRecordIterator struct {
	rows *sql.Rows
}

func (r IntPeriodRecordIterator) Next() (warehouse.PeriodRecord[int64], bool) {
	r.rows.Next()

	var value warehouse.PeriodRecord[int64]
	var rawStart int64
	var rawEnd int64
	err := r.rows.Scan(&value.Value, &rawStart, &rawEnd)
	if err != nil {
		return value, false
	}

	value.Start = time.Unix(rawStart, 0)
	value.End = time.Unix(rawEnd, 0)

	return value, true
}

func (r IntPeriodRecordIterator) Err() error {
	return r.rows.Err()
}

func (r IntPeriodRecordIterator) Close() error {
	return r.rows.Close()
}

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
	_, err := Exec(r.db, "INSERT INTO realrecords(value, timestamp, bucket) values(?, ?, ?)", value.Value, value.Timestamp.Unix(), bucket)
	return err
}

func (r *Real) Read(bucket string, start, end time.Time) ([]warehouse.Record[float64], error) {
	rows, err := Query(r.db, "SELECT value, timestamp FROM realrecords WHERE timestamp >= ? AND timestamp <= ? AND bucket = ?", start.Unix(), end.Unix(), bucket)
	if err != nil {
		return nil, err
	}

	var values []warehouse.Record[float64]

	for rows.Next() {
		var value warehouse.Record[float64]

		var rawTimestamp int64
		err = rows.Scan(&value.Value, &rawTimestamp)
		if err != nil {
			return nil, err
		}

		value.Timestamp = time.Unix(rawTimestamp, 0)

		values = append(values, value)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

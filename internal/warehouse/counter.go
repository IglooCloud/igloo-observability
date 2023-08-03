package warehouse

import (
	"time"
)

type Counter struct {
	Storage PeriodStore[int64]
}

// Write a record to the database
func (c *Counter) Write(bucket string, value PeriodRecord[int64]) error {
	return c.Storage.Write(bucket, value)
}

// Read records in the time range
func (c *Counter) Read(bucket string, start time.Time, end time.Time) ([]PeriodRecord[int64], error) {
	records := make([]PeriodRecord[int64], 0)
	recordIterator, err := c.Storage.Read(bucket, start, end)
	if err != nil {
		return nil, err
	}

	for record, ok := recordIterator.Next(); ok; record, ok = recordIterator.Next() {
		records = append(records, record)
	}
	if err := recordIterator.Err(); err != nil {
		return nil, err
	}
	recordIterator.Close()

	return records, nil
}

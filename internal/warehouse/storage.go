package warehouse

import (
	"errors"
	"time"
)

var NO_RECORD_FOUND = errors.New("no records found for this time range")

type Record[T any] struct {
	Value     T
	Timestamp time.Time
}

type Iterator[T any] interface {
	Next() (T, bool)
	Err() error
	Close() error
}

type RecordIterator[T any] Iterator[Record[T]]

type Store[T any] interface {
	Write(bucket string, value Record[T]) error
	Read(bucket string, start time.Time, end time.Time) (RecordIterator[T], error)
	ReadValues(bucket string, start time.Time, end time.Time) (Iterator[T], error)
}

type PeriodRecord[T any] struct {
	Value T
	Start time.Time
	End   time.Time
}

type PeriodRecordIterator[T any] Iterator[PeriodRecord[T]]

type PeriodStore[T any] interface {
	Write(bucket string, value PeriodRecord[T]) error
	Read(bucket string, start time.Time, end time.Time) (PeriodRecordIterator[T], error)
}

package warehouse

import "time"

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

package warehouse

import "time"

type Record[T any] struct {
	Value     T
	Timestamp time.Time
}

type Store[T any] interface {
	Write(bucket string, value Record[T]) error
	Read(bucket string, start time.Time, end time.Time) ([]Record[T], error)
}

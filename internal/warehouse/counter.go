package warehouse

import "time"

type Counter struct {
	Storage Store[float64]
}

func (c *Counter) Write(bucket string, value Record[float64]) error {
	return c.Storage.Write(bucket, value)
}

func (c *Counter) Read(bucket string, start time.Time, end time.Time) ([]Record[float64], error) {
	return c.Storage.Read(bucket, start, end)
}

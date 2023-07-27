package warehouse

import "time"

type Gauge struct {
	Storage Store[float64]
}

func (g *Gauge) Write(bucket string, value Record[float64]) error {
	return g.Storage.Write(bucket, value)
}

func (g *Gauge) Read(bucket string, start time.Time, end time.Time) ([]Record[float64], error) {
	return g.Storage.Read(bucket, start, end)
}

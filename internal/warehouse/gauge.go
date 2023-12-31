package warehouse

import (
	"sort"
	"time"

	stat "gonum.org/v1/gonum/stat"
)

type Gauge struct {
	Storage Store[float64]
}

// Write a record to the database
func (g *Gauge) Write(bucket string, value Record[float64]) error {
	return g.Storage.Write(bucket, value)
}

// Read records in the time range
func (g *Gauge) Read(bucket string, start time.Time, end time.Time) ([]Record[float64], error) {
	records := make([]Record[float64], 0)
	recordIterator, err := g.Storage.Read(bucket, start, end)
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

// ReadValues reads the values associated to the records in the time range
func (g *Gauge) ReadValues(bucket string, start time.Time, end time.Time) ([]float64, error) {
	records := make([]float64, 0)
	iterator, err := g.Storage.ReadValues(bucket, start, end)
	if err != nil {
		return nil, err
	}

	for record, ok := iterator.Next(); ok; record, ok = iterator.Next() {
		records = append(records, record)
	}
	if err := iterator.Err(); err != nil {
		return nil, err
	}
	iterator.Close()

	return records, nil
}

// Mean returns the mean value of the records in the time range
func (g *Gauge) Mean(bucket string, start time.Time, end time.Time) (float64, error) {
	iterator, err := g.Storage.ReadValues(bucket, start, end)
	if err != nil {
		return 0, err
	}

	var sum float64 = 0
	var count int = 0
	for value, ok := iterator.Next(); ok; value, ok = iterator.Next() {
		sum += value
		count++
	}
	iterator.Close()

	if err := iterator.Err(); err != nil {
		return 0, err
	} else if count == 0 {
		return 0, NO_RECORD_FOUND
	}

	return sum / float64(count), nil
}

// Min returns the minimum value of the records in the time range
func (g *Gauge) Min(bucket string, start time.Time, end time.Time) (float64, error) {
	iterator, err := g.Storage.ReadValues(bucket, start, end)
	if err != nil {
		return 0, err
	}

	firstValue, ok := iterator.Next()
	if !ok {
		return 0, NO_RECORD_FOUND
	}

	var min float64 = firstValue
	for value, ok := iterator.Next(); ok; value, ok = iterator.Next() {
		if value < min {
			min = value
		}
	}
	iterator.Close()

	if err := iterator.Err(); err != nil {
		return 0, err
	}

	return min, nil
}

// Max returns the maximum value of the records in the time range
func (g *Gauge) Max(bucket string, start time.Time, end time.Time) (float64, error) {
	iterator, err := g.Storage.ReadValues(bucket, start, end)
	if err != nil {
		return 0, err
	}

	firstValue, ok := iterator.Next()
	if !ok {
		return 0, NO_RECORD_FOUND
	}

	var max float64 = firstValue
	for value, ok := iterator.Next(); ok; value, ok = iterator.Next() {
		if value > max {
			max = value
		}
	}
	iterator.Close()

	if err := iterator.Err(); err != nil {
		return 0, err
	}

	return max, nil
}

// Percentile returns the percentile value of the records in the time range
func (g *Gauge) Percentile(bucket string, start time.Time, end time.Time, percentile float64) (float64, error) {
	values, err := g.ReadValues(bucket, start, end)
	if err != nil {
		return 0, err
	} else if len(values) == 0 {
		return 0, NO_RECORD_FOUND
	}

	sort.Float64s(values)

	return stat.Quantile(percentile, stat.LinInterp, values, nil), nil
}

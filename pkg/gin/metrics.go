package observability

import (
	"sync"
	"time"
)

// counter is an interface that allows to atomically
// update and read an integer value.
type counter interface {
	Add(delta int64) (new int64)
	CompareAndSwap(old int64, new int64) (swapped bool)
	Load() int64
	Store(val int64)
	Swap(new int64) (old int64)
}

type metrics struct {
	lock          *sync.RWMutex
	StartTime     time.Time
	TotalRequests counter
	Error400      counter
	Error500      counter
	TotalLatency  counter
}

func (m *metrics) Record(latency time.Duration, status int) {
	// We can use a read lock here because we are only
	// doing atomic operations on the counters.
	m.lock.RLock()
	defer m.lock.RUnlock()

	m.TotalRequests.Add(1)
	m.TotalLatency.Add(latency.Microseconds())

	switch {
	case status >= 500 && status < 600:
		m.Error500.Add(1)
	case status >= 400 && status < 500:
		m.Error400.Add(1)
	}
}

type OutputMetrics struct {
	TotalRequests int64
	Error400      int64
	Error500      int64
	TotalLatency  int64
	StartTime     time.Time
	EndTime       time.Time
}

func (m *metrics) Reset() OutputMetrics {
	m.lock.Lock()
	defer m.lock.Unlock()
	endTime := time.Now()

	totalRequests := m.TotalRequests.Swap(0)
	error400 := m.Error400.Swap(0)
	error500 := m.Error500.Swap(0)
	totalLatency := m.TotalLatency.Swap(0)

	startTime := m.StartTime
	m.StartTime = endTime

	return OutputMetrics{
		TotalRequests: totalRequests,
		Error400:      error400,
		Error500:      error500,
		TotalLatency:  totalLatency,
		StartTime:     startTime,
		EndTime:       endTime,
	}
}

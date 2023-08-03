package observability

import (
	"sync/atomic"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/collector"
	"github.com/gin-gonic/gin"
)

func New() (middleware gin.HandlerFunc, routeHandler gin.HandlerFunc) {
	m := metrics{
		StartTime:     time.Now(),
		TotalRequests: &atomic.Int64{},
		TotalLatency:  &atomic.Int64{},
		Error400:      &atomic.Int64{},
		Error500:      &atomic.Int64{},
	}

	return newMiddleware(m), newRouteHandler(m)
}

func newMiddleware(m metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		m.Record(time.Since(startTime), c.Writer.Status())
	}
}

func newRouteHandler(m metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		output := m.Reset()

		payload := collector.FetchResponse{
			Gauges: map[string]collector.GaugeRecord{},
			Counters: map[string]collector.CounterRecord{
				"requests": {
					Value: output.TotalRequests,
					Start: output.StartTime,
					End:   output.EndTime,
				},
				"totalLatency": {
					Value: output.TotalLatency,
					Start: output.StartTime,
					End:   output.EndTime,
				},
				"error400": {
					Value: output.Error400,
					Start: output.StartTime,
					End:   output.EndTime,
				},
				"error500": {
					Value: output.Error500,
					Start: output.StartTime,
					End:   output.EndTime,
				},
			},
		}

		c.JSON(200, payload)
	}
}
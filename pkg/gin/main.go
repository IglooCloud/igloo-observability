// This package provides a Gin middleware and route handler for collecting metrics
// from Gin web applications and exposing them in a igloo-observability compatible format.
//
// To instrument a Gin application, add the middleware to the Gin router and add the route handler:
//
//	observabilityMiddleware, metricsHandler := observability.New()
//	r.Use(observabilityMiddleware)
//	r.GET("/metrics", metricsHandler)
//
// To secure the metrics endpoint, set the IGLOO_OBSERVABILITY_SECRET environment variable.
package observability

import (
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/collector"
	"github.com/gin-gonic/gin"
)

// New returns a Gin middleware and route handler for collecting metrics from Gin web applications
func New() (middleware gin.HandlerFunc, routeHandler gin.HandlerFunc) {
	m := metrics{
		lock:          &sync.RWMutex{},
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
	secret := os.Getenv("IGLOO_OBSERVABILITY_SECRET")

	return func(c *gin.Context) {
		if secret != "" && c.GetHeader("igloo-observability-secret") != secret {
			c.String(http.StatusUnauthorized, "Wrong secret")
			return
		}

		output := m.Reset()

		payload := collector.FetchResponse{
			Gauges: map[string]collector.GaugeRecord{},
			Counters: map[string]collector.CounterRecord{
				"http.requests": {
					Value: output.TotalRequests,
					Start: output.StartTime,
					End:   output.EndTime,
				},
				"http.totalLatency": {
					Value: output.TotalLatency,
					Start: output.StartTime,
					End:   output.EndTime,
				},
				"http.error400": {
					Value: output.Error400,
					Start: output.StartTime,
					End:   output.EndTime,
				},
				"http.error500": {
					Value: output.Error500,
					Start: output.StartTime,
					End:   output.EndTime,
				},
			},
		}

		c.JSON(200, payload)
	}
}

package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/warehouse"
	"github.com/gin-gonic/gin"
)

func recordsToCSV(records []warehouse.Record[float64]) string {
	var csv string
	for _, record := range records {
		csv += fmt.Sprintf("%s,%f\n", record.Timestamp.Format(time.DateTime), record.Value)
	}
	return csv
}

func registerGaugeRoutes(r *gin.Engine, gauge warehouse.Gauge) {
	r.GET("/gauge/:name", func(c *gin.Context) {
		name := c.Param("name")
		var query struct {
			Start  time.Time `form:"start" time_format:"2006-01-02T15:04:05Z"`
			End    time.Time `form:"end" time_format:"2006-01-02T15:04:05Z"`
			Format string    `form:"format"`
		}
		if err := c.ShouldBindQuery(&query); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		if query.Start.IsZero() {
			query.Start = time.Now().Add(-3 * time.Hour)
		}
		if query.End.IsZero() {
			query.End = time.Now()
		}
		if query.Format == "" {
			query.Format = "json"
		}

		series, err := gauge.Read(name, query.Start, query.End)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		switch query.Format {
		case "json":
			c.JSON(http.StatusOK, series)
		case "csv":
			c.String(http.StatusOK, recordsToCSV(series))
		default:
			c.String(http.StatusBadRequest, "invalid format")
		}
	})
}

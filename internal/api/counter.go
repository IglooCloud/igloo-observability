package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/warehouse"
	"github.com/gin-gonic/gin"
)

func counterRecordsToCSV(records []warehouse.PeriodRecord[int64]) string {
	var csv string
	for _, record := range records {
		csv += fmt.Sprintf("%s,%s,%d\n", record.Start.Format(time.DateTime), record.End.Format(time.DateTime), record.Value)
	}
	return csv
}

func registerCounterRoutes(r *gin.Engine, counter warehouse.Counter) {
	r.GET("/counter/:name", func(c *gin.Context) {
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

		series, err := counter.Read(name, query.Start, query.End)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		switch query.Format {
		case "json":
			c.JSON(http.StatusOK, series)
		case "csv":
			c.String(http.StatusOK, counterRecordsToCSV(series))
		default:
			c.String(http.StatusBadRequest, "invalid format")
		}
	})
}

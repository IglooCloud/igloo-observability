package collector

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/warehouse"
)

type Endpoint struct {
	URL    string
	Secret string
	Bucket string
}

type GaugeRecord struct {
	Value     float64   `json:"v"`
	Timestamp time.Time `json:"t"`
}
type FetchResponse struct {
	Gauges map[string]GaugeRecord `json:"g"`
}

// Fetch data from endpoint and write to gauge
func Fetch(endpoint Endpoint, gauge warehouse.Gauge) error {
	// Fetch data from endpoint
	req, err := http.NewRequest("GET", endpoint.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("igloo-observability-secret", endpoint.Secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse payload
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var payload FetchResponse
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return err
	}

	return writeResponse(endpoint.Bucket, payload, gauge)
}

func writeResponse(bucket string, resp FetchResponse, gauge warehouse.Gauge) error {
	for name, record := range resp.Gauges {
		err := gauge.Write(bucket+"."+name, warehouse.Record[float64](record))
		if err != nil {
			return err
		}
	}
	return nil
}

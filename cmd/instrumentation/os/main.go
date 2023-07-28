package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/collector"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	port := os.Getenv("PORT")
	secret := os.Getenv("SECRET")
	if port == "" {
		port = "6001"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This is an OS instrumentation server, see github.com/IglooCloud/igloo-observability for more information.")
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if secret != "" && r.Header.Get("igloo-observability-secret") != secret {
			http.Error(w, "Wrong secret", http.StatusUnauthorized)
			return
		}

		payload := collector.FetchResponse{
			Gauges: map[string]collector.GaugeRecord{},
		}

		// Load metrics
		err := LoadMemoryMetrics(&payload)
		if err != nil {
			http.Error(w, "Failed to load memory metrics", http.StatusInternalServerError)
			return
		}
		err = LoadCPUMetrics(&payload)
		if err != nil {
			http.Error(w, "Failed to load cpu metrics", http.StatusInternalServerError)
			return
		}
		err = LoadDiskMetrics(&payload)
		if err != nil {
			http.Error(w, "Failed to load disk metrics", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	fmt.Printf("Server listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

const B_to_GB float64 = 1.0 / (1024.0 * 1024.0 * 1024.0)
const B_to_gB float64 = 1.0 / (1000.0 * 1000.0 * 1000.0)

func LoadMemoryMetrics(payload *collector.FetchResponse) error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	payload.Gauges["memory.total"] = collector.GaugeRecord{
		Value:     float64(memory.Total) * B_to_GB,
		Timestamp: time.Now(),
	}
	payload.Gauges["memory.used"] = collector.GaugeRecord{
		Value:     float64(memory.Used) * B_to_GB,
		Timestamp: time.Now(),
	}

	return nil
}

func LoadCPUMetrics(payload *collector.FetchResponse) error {
	loadAverage, err := load.Avg()
	if err != nil {
		return err
	}
	payload.Gauges["cpu.load1"] = collector.GaugeRecord{
		Value:     float64(loadAverage.Load1),
		Timestamp: time.Now(),
	}

	return nil
}

func LoadDiskMetrics(payload *collector.FetchResponse) error {
	disk, err := disk.Usage("/")
	if err != nil {
		return err
	}
	payload.Gauges["disk.total"] = collector.GaugeRecord{
		Value:     float64(disk.Total) * B_to_gB,
		Timestamp: time.Now(),
	}
	payload.Gauges["disk.used"] = collector.GaugeRecord{
		Value:     float64(disk.Used) * B_to_gB,
		Timestamp: time.Now(),
	}

	return nil
}

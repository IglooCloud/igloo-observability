package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/IglooCloud/igloo-observability/internal/api"
	"github.com/IglooCloud/igloo-observability/internal/collector"
	"github.com/IglooCloud/igloo-observability/internal/config"
	"github.com/IglooCloud/igloo-observability/internal/log"
	"github.com/IglooCloud/igloo-observability/internal/storage"
	"github.com/IglooCloud/igloo-observability/internal/warehouse"
)

func fetchWorker(endpointStream chan collector.Endpoint, gauge warehouse.Gauge) {
	var logger = log.Default().Service("worker")

	for endpoint := range endpointStream {
		err := collector.Fetch(endpoint, gauge)
		if err != nil {
			logger.Error(err)
		}
	}
}

func main() {
	config, err := config.LoadTOML("./example.toml")
	if err != nil {
		panic(err)
	}

	db, err := storage.Connect(config.DBPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	real := storage.RealFromDB(db)
	gauge := warehouse.Gauge{Storage: real}

	endpointStream := collector.RequestStream(config.Collector)
	go fetchWorker(endpointStream, gauge)

	go api.Start(gauge, config.API)

	// Block until a signal is received
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

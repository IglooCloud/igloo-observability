package main

import (
	"fmt"
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

func fetchWorker(endpointStream chan collector.Endpoint, gauge warehouse.Gauge, counter warehouse.Counter) {
	var logger = log.Default().Service("worker")

	for endpoint := range endpointStream {
		err := collector.Fetch(endpoint, gauge, counter)
		if err != nil {
			logger.Error(err)
		}
	}
}

func main() {
	// Load config
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: server <config-path>")
		os.Exit(1)
	}

	configPath := os.Args[1]
	config, err := config.LoadTOML(configPath)
	if err != nil {
		panic(err)
	}

	// Connect to database
	db, err := storage.Connect(config.DBPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Inject dependencies
	real := storage.RealFromDB(db)
	periodInt := storage.PeriodIntFromDB(db)
	gauge := warehouse.Gauge{Storage: real}
	counter := warehouse.Counter{Storage: periodInt}

	endpointStream := collector.RequestStream(config.Collector)
	go fetchWorker(endpointStream, gauge, counter)

	go api.Start(gauge, counter, config.API)

	// Block until a signal is received
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

package main

import (
	"github.com/IglooCloud/igloo-observability/internal/collector"
	"github.com/IglooCloud/igloo-observability/internal/config"
	"github.com/IglooCloud/igloo-observability/internal/storage"
	"github.com/IglooCloud/igloo-observability/internal/warehouse"
)

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
	for endpoint := range endpointStream {
		err := collector.Fetch(endpoint, gauge)
		if err != nil {
			panic(err)
		}
	}
}

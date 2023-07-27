package main

import (
	"fmt"
	"time"

	"github.com/IglooCloud/igloo-observability/internal/storage"
	"github.com/IglooCloud/igloo-observability/internal/warehouse"
)

func main() {
	db, err := storage.Connect("./foo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	real := storage.RealFromDB(db)
	gauge := warehouse.Gauge{Storage: real}

	err = gauge.Write("test", warehouse.Record[float64]{Value: 1.0, Timestamp: time.Now()})
	if err != nil {
		panic(err)
	}
	err = gauge.Write("test", warehouse.Record[float64]{Value: 2.0, Timestamp: time.Now()})
	if err != nil {
		panic(err)
	}

	values, err := real.Read("test", time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(values)
}

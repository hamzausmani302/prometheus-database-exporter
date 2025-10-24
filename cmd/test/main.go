package main

import (
	"fmt"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	cfg := config.GetConfig("example", logger)
	fmt.Println(cfg.EnableApi)
	fmt.Println(cfg.EnableCollector)
	fmt.Println(cfg.Store.StoreType)
	
}
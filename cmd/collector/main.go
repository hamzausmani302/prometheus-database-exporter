package main

import (
	"fmt"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
)

/*
entry point for collector which collects metrics from the database
and puts them in store from where prmetheus can scrape them via API.
*/
func main() {
	fmt.Println("Collector started")	
	cfg := config.GetConfig("example")
	fmt.Println(cfg)
}
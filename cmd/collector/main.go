package main

import (
	"fmt"

	"github.com/hamzausmani302/prometheus-database-exporter/config"

	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"

	"github.com/hamzausmani302/prometheus-database-exporter/internal/factories"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/utils"

	"github.com/sirupsen/logrus"
)

/*
entry point for collector which collects metrics from the database
and puts them in store from where prmetheus can scrape them via API.
*/
func main() {

	logger := logrus.New()
	fmt.Println("Collector started")	
	
	// Read config from file

	cfg := config.GetConfig("example", logger)
	logger.Debug(cfg)

	//Creating Datasource
	ds := factories.NewDatasourceFactory(logger, &cfg).Create(cfg.DataSource[0])
	ds.Connect()
	df := ds.GetData(datasource.SQLQuery{
		Query: "select well_name, well_id from t_well;",
	})	
	fmt.Println(df)
	ds.Close()

	// Creating Store
	c := factories.NewCacheStoreFactory(logger, &cfg).Create(cfg.Store)
	b , _ := utils.FataFrameToCSVBytes(df)
	c.Set("key1", b, 5)
	stop := false
	for  stop == false{
		data , err := c.Get("key1")
		if err != nil {
			fmt.Println(err)
		}
		d := utils.DataFrameFromCSVBytes(data)
		fmt.Println(d)

	}
	// Creating schduler(Inject dependencies)
	

	// Execute schduler

	
	// On interrupt, close datasource, and all the connections
	


}
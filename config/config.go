package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

/*
Contains all the configuration settings for the application and exporter.
Defines structs for Store, Collector, DataSource and main Application config.
Also defines enums for DataSourceType and CollectorType.
*/
// Enums for DataSourceType and CollectorType
type DataSourceType string

// Enums for DataSourceType and CollectorType
type CollectorType string

// Enums for Schedulertype
type SchedulerType string

const (
	// Enum mapping for DataSourceType
	SQL DataSourceType = "SQL"
	// Enum mapping for CollectorType
	Prometheus CollectorType = "Prometheus"
	// enum mapping for SchedulerType
	Memory SchedulerType = "memory"
	Sqlite SchedulerType = "sqlite"
	Redis  SchedulerType = "redis"
)

/*
configuration for the task schduler
*/
type SchedulerConfig struct {
	Storage  SchedulerType           `yaml:"storage"`
	Metadata SchedulerMetadataConfig `yaml:"metadata"`
}
type SchedulerMetadataConfig struct {
	ConnectionDetails map[string]string `yaml:"connectionDetails"`
}

/*
Configuration structs for the Store.
*/
type StoreConfig struct {
	// Type of the store enum (InMemory, Redis)
	StoreType string `yaml:"type"`
	// Metadata for the store (Specifying connection details)
	Metadata StoreConfigMetadataConfig `yaml:"metadata"`
}
type StoreConfigMetadataConfig struct {
	ConnectionDetails map[string]string `yaml:"connectionDetails"`
}

/*
Configuration structs for the Collector.
*/
type CollectorConfig struct {
	// Type of collector enum (Prometheus)
	CollectType CollectorType `yaml:"type"`
	// Metadata for the collector (Specifying additional details)
	// Makes it easy to implement new features without chaning the config much
	Metadata map[string]string `yaml:"metadata"`
}

/*
Configuration structs for the DataSource.
*/
type DataSourceMetadataConfig struct {
	// Connection details for the datasource (host, port, username, password, dbname etc)
	ConnectionDetails map[string]string `yaml:"connectionDetails"`
}

/*
Configuration structs for the DataSource.
*/
type DataSourceConfig struct {
	// Name of the datasource (e.g. Postgres, MySQL), will be used in metrics
	Name string `yaml:"name"`
	// Type of the datasource enum (SQL)
	DataSourceType DataSourceType `yaml:"type"`
	// Metadata for the datasource (Specifying connection details)
	Metadata DataSourceMetadataConfig `yaml:"metadata"`
}

/*
Main configuration struct for the application.
Containing all the sub-configurations.
*/
type ApplicationConfig struct {
	// Configuration for the Schduler
	Scheduler SchedulerConfig `yaml:"schedulerConfig"`
	// Configuration for the Store
	Store StoreConfig `yaml:"storeConfig"`
	// Configuration for the Collector
	Collector CollectorConfig `yaml:"collectorConfig"`
	// Configuration for the DataSource
	DataSource []DataSourceConfig `yaml:"dataSourceConfig"`
	// Queries to be executed to fetch metrics
	Queries []map[string]interface{} `yaml:"queries"`
}

func (cfg *ApplicationConfig) readConfigData(data []byte) {
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}
}

var appCfg *ApplicationConfig

func GetConfig(env string, logger *logrus.Logger) ApplicationConfig {
	var applicationConfig ApplicationConfig
	if appCfg == nil {
		logger.SetLevel(logrus.DebugLevel)
		configFilePath := path.Join("config", fmt.Sprintf("config.%s.yaml", env))
		logger.Infof("Reading config file: %s ", configFilePath)
		// Read the config file

		content, err := os.ReadFile(configFilePath)
		if err != nil {
			logger.Fatalf("Error reading file: %v", err)
			panic("There is a problem reading the file...")
		}
		applicationConfig.readConfigData(content)
		// fmt.Println(applicationConfig)
		appCfg = &applicationConfig
	}
	return *appCfg
}

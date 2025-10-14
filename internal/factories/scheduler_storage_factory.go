package factories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/algorythma/go-scheduler/storage"
	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/sirupsen/logrus"
)

/*
Factory class to initiate the configurations and provide abstraction over the creation of
datasource objects
*/
type SchedulerStorageFactory struct {
	logger *logrus.Logger
	cfg    *config.ApplicationConfig
}

func (dsf *SchedulerStorageFactory) Create(schedulerConfig config.SchedulerConfig) (storage.TaskStore, error) {
	if strings.EqualFold(strings.ToLower(string(schedulerConfig.Storage)) , strings.ToLower(string(config.Memory)) ){
		return storage.NewMemoryStorage(), nil
	} else if strings.EqualFold(strings.ToLower(string(schedulerConfig.Storage)), strings.ToLower(string(config.Sqlite))) {
		if schedulerConfig.Metadata.ConnectionDetails["dbName"] == "" {
			dsf.logger.Warn("dbName not provided")
			return nil, errors.New("DbName not provided")
		}
		dsf.logger.Debugf("sc = %s", schedulerConfig)
		strg := storage.NewSqlite3Storage(storage.Sqlite3Config{
			DbName: "test123",
		})
		if err := strg.Connect(); err != nil {
			return nil, err
		}
		if err := strg.Initialize(); err != nil {
			return nil, err
		}
		return strg, nil
	} else if strings.EqualFold(strings.ToLower(string(schedulerConfig.Storage)), strings.ToLower(string(config.Redis))) {
		if schedulerConfig.Metadata.ConnectionDetails["dbName"] == "" {
			dsf.logger.Warn("dbName not provided")
			return nil, errors.New("DbName not provided")
		}
		dsf.logger.Debug("sc", schedulerConfig)
		strg, _ := storage.NewRedisStorage(storage.RedisConfig{
			Host:     "test",
			Port:     6379,
			Password: "",
			Db:       0,
		})

		return strg, nil
	}
	dsf.logger.Fatalf("Invalid Storage type: %s", schedulerConfig.Storage)
	return nil, fmt.Errorf("Invalid Storage type: %s", schedulerConfig.Storage)
}

func NewSchdulerStorageFactory(logger *logrus.Logger, cfg *config.ApplicationConfig) *SchedulerStorageFactory {
	return &SchedulerStorageFactory{logger, cfg}
}

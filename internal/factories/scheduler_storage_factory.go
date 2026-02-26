package factories

import (
	"errors"
	"fmt"
	"strconv"
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
	if strings.EqualFold(string(schedulerConfig.Storage), string(config.Memory)) {
		return storage.NewMemoryStorage(), nil
	} else if strings.EqualFold(string(schedulerConfig.Storage), string(config.Sqlite)) {
		dbName := schedulerConfig.Metadata.ConnectionDetails["dbName"]
		if dbName == "" {
			dsf.logger.Warn("dbName not provided")
			return nil, errors.New("dbName not provided")
		}
		dsf.logger.Debugf("creating sqlite storage with dbName=%s", dbName)
		strg := storage.NewSqlite3Storage(storage.Sqlite3Config{
			DbName: dbName,
		})
		if err := strg.Connect(); err != nil {
			return nil, err
		}
		if err := strg.Initialize(); err != nil {
			return nil, err
		}
		return strg, nil
	} else if strings.EqualFold(string(schedulerConfig.Storage), string(config.Redis)) {
		host := schedulerConfig.Metadata.ConnectionDetails["host"]
		if host == "" {
			return nil, errors.New("host not provided for redis scheduler storage")
		}
		port, err := strconv.Atoi(schedulerConfig.Metadata.ConnectionDetails["port"])
		if err != nil {
			port = 6379
		}
		db, err := strconv.Atoi(schedulerConfig.Metadata.ConnectionDetails["db"])
		if err != nil {
			db = 0
		}
		dsf.logger.Debugf("creating redis scheduler storage host=%s port=%d", host, port)
		strg, err := storage.NewRedisStorage(storage.RedisConfig{
			Host:     host,
			Port:     port,
			Password: schedulerConfig.Metadata.ConnectionDetails["password"],
			Db:       db,
		})
		if err != nil {
			return nil, err
		}
		return strg, nil
	}
	return nil, fmt.Errorf("invalid storage type: %s", schedulerConfig.Storage)
}

func NewSchdulerStorageFactory(logger *logrus.Logger, cfg *config.ApplicationConfig) *SchedulerStorageFactory {
	return &SchedulerStorageFactory{logger, cfg}
}

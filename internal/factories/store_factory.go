package factories

import (
	"fmt"
	"strconv"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/cache"
	"github.com/sirupsen/logrus"
)

/*
Factory class to initiate the configurations and provide abstraction over the creation of
datasource objects
*/
type CacheStoreFactory struct {
	logger *logrus.Logger
	cfg    *config.ApplicationConfig
}

func (dsf *CacheStoreFactory) Create(storeConfig config.StoreConfig) (cache.ICache, error) {
	dsf.logger.Debugf("Creating %s store", storeConfig.StoreType)

	if storeConfig.StoreType == "local" {
		return cache.NewLocaltimeCache(), nil
	} else if storeConfig.StoreType == "redis" {
		port, err := strconv.Atoi(storeConfig.Metadata.ConnectionDetails["port"])
		if err != nil {
			dsf.logger.Warn("Store Config Port not specified, defaulting to 6379")
			port = 6379
		}
		c := cache.NewRedisCache(cache.RedisConnectionSettings{
			Host:     storeConfig.Metadata.ConnectionDetails["host"],
			Port:     port,
			Password: storeConfig.Metadata.ConnectionDetails["password"],
		})
		if c == nil {
			return nil, fmt.Errorf("failed to connect to redis at %s:%d", storeConfig.Metadata.ConnectionDetails["host"], port)
		}
		return c, nil
	}
	return nil, fmt.Errorf("invalid store type: %s", storeConfig.StoreType)
}

func NewCacheStoreFactory(logger *logrus.Logger, cfg *config.ApplicationConfig) *CacheStoreFactory {
	return &CacheStoreFactory{logger: logger, cfg: cfg}
}

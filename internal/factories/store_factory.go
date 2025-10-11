package factories

import (
	"strconv"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/cache"
	"github.com/sirupsen/logrus"
)

/*
Factory class to initiate the configurations and provide abstraction over the creation of
datasource objects
*/
type CacheStoreFactory struct{
	logger *logrus.Logger
	cfg *config.ApplicationConfig
}

func (dsf *CacheStoreFactory) Create( storeConfig config.StoreConfig) cache.ICache {
	dsf.logger.Debugf("storeConfig", storeConfig)
	dsf.logger.Debugf("Creating %s store", storeConfig.StoreType)
	
	if storeConfig.StoreType == "local"{
		return cache.NewLocaltimeCache()
	}else if storeConfig.StoreType == "redis"{
		port, err := strconv.Atoi(storeConfig.Metadata.ConnectionDetails["port"])
		if err != nil{
			dsf.logger.Warn("Store Config Port not specified", err)
		}
		return cache.NewRedisCache(cache.RedisConnectionSettings{
			Host: storeConfig.Metadata.ConnectionDetails["host"],
			Port: port,
			Password: storeConfig.Metadata.ConnectionDetails["password"],
		})
	}
	dsf.logger.Fatalf("Invalid Store provided : %s", storeConfig.StoreType)
	return nil
}

func NewCacheStoreFactory(logger *logrus.Logger, cfg *config.ApplicationConfig) *CacheStoreFactory{
	return &CacheStoreFactory{logger: logger, cfg: cfg};
} 
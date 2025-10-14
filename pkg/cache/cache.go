package cache

import (
	"errors"
	"time"

	_cache "github.com/aleksiumish/in-memory-cache"
)

// Interface for the cache to implement
type ICache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte, expiresIn int64) error
	GetCacheType() string
}

// Implementation of time based cache allowing to ttl feature
type LocalTimeCache struct {
	Cache     _cache.Cache
	TtlLookup map[string]int64
}

func (ltc *LocalTimeCache) GetCacheType() string {
	return "localcache"
}

// Get value for specified key
func (ltc *LocalTimeCache) Get(key string) ([]byte, error) {
	var value []byte
	gotValue := ltc.Cache.Get(key)
	if gotValue != nil {
		value = gotValue.([]byte)
		now := time.Now()
		if value != nil && ltc.TtlLookup[key] > now.UnixMilli() {
			return value, nil
		}
	}

	return nil, errors.New("Not Found")
}

/*
	Set the key and ttl value
	Input:
		key: key of the item
		data: value for the item
		expiresIn: duration in seconds for the item to expire
*/
func (ltc *LocalTimeCache) Set(key string, data []byte, expiresIn int64) error {
	ltc.Cache.Set(key, data)
	now := time.Now()

	ltc.TtlLookup[key] = now.UnixMilli() + (expiresIn * 1000)
	return nil
}

// Create an instance of
func NewLocaltimeCache() ICache {
	return &LocalTimeCache{Cache: *_cache.NewCache(), TtlLookup: map[string]int64{}}
}

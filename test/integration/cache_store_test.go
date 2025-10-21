//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/hamzausmani302/prometheus-database-exporter/pkg/cache"
)

// Test simple memory caching Store
// Test Redis connectivity
// Test SQL database connnectivity
//
func TestFunctionStore(t *testing.T) {
	cacheStore := cache.NewRedisCache(cache.RedisConnectionSettings{
		Host: "localhost",
		Port: 6379,
		Password: "",
	})
	expected := "value"
	key := "testkey" 
	if cacheStore == nil{
		t.Errorf("unable to establish connection to Redis")
		return
	}
	if err := cacheStore.Set(key, []byte(expected), 10); err != nil {
		t.Error("Cannot set the value to redis")
	}
	var data string
	if value, err := cacheStore.Get(key); err != nil{
		t.Error(err)
	}else{
		data= string(value)
	}

	if data != expected{
		t.Errorf("Expected = %s, Got = %s", expected, data)
	}


	
}
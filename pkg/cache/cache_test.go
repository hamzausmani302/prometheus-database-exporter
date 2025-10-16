package cache

import (
	"testing"
	"time"
)

// check if the instance is created properly and interface methods are working
func TestCacheInstance(t *testing.T) {	
	var cache ICache = NewLocaltimeCache()
	if cache.GetCacheType() != "localcache" {
		t.Errorf("Expected cache type localcache, got %s", cache.GetCacheType())
	}
}


// check the implementation is working fine for localtime cache
func TestLocaltimeCache(t *testing.T) {
	cache := NewLocaltimeCache()
	key := "testKey"
	value := []byte("testValue")
	expiry := int64(5) // 5 seconds
	if err1 := cache.Set(key, value, expiry); err1 != nil {
		t.Errorf("Error setting value in cache: %s", err1.Error())
	}
	timer := time.NewTimer(6 * time.Second)
	if val, err := cache.Get(key); err != nil || string(val) != string(value) {
		t.Errorf("Expected value %s, got %s, err: %v", string(value), string(val), err)
	}
	<-timer.C
	if val ,err := cache.Get(key); err == nil {
		t.Errorf("Expected no value , got = %s", string(val))
	}
}
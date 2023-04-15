package cache

import (
	"log"
	"time"
)

// NewCoreCache initializes a new core cache and also launch a new goroutine for running expireCheck in the background
func NewCoreCache(expireInterval time.Duration) *CoreCache {
	cc := &CoreCache{
		cacheDataMap:        make(map[string]CacheData),
		gracefulStopChannel: make(chan struct{}),
	}

	cc.wg.Add(1)
	go func(expireInterval time.Duration) {
		defer cc.wg.Done()
		cc.expireCheck(expireInterval)
	}(expireInterval)

	return cc

}

// expireCheck is the process that runs in the background to check expiration of a key and also to check gracefulStop signal
func (cc *CoreCache) expireCheck(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-cc.gracefulStopChannel:
			//delete all keys from cache when doing graceful stop
			cc.rwmu.Lock()
			log.Println("[shutting down] begin cleaning up cache")
			for key, _ := range cc.cacheDataMap {
				log.Println("[shutting down] deleting key", key)
				delete(cc.cacheDataMap, key)
			}
			log.Println("[shutting down] finished cleaning up cache")
			cc.rwmu.Unlock()
			return
		case <-t.C:
			cc.rwmu.Lock()
			for key, data := range cc.cacheDataMap {
				if data.expirationEpochTimestamp <= time.Now().Unix() {
					log.Println("[expiration check] deleting key", key)
					delete(cc.cacheDataMap, key)
				}
			}
			cc.rwmu.Unlock()
		}
	}
}

// StopExpireCheck is the function to be called during shutting down the application (for graceful shutdown) to delete all the keys in the cache.
func (cc *CoreCache) StopExpireCheck() {
	close(cc.gracefulStopChannel)
	cc.wg.Wait()
}

// Set is used to set a new (or override an existing) key along with its value into the cache. The default expiration time for each key is 30 Minutes(1800 seconds)
func (cc *CoreCache) Set(key, value string) {
	cc.rwmu.Lock()
	defer cc.rwmu.Unlock()

	cc.cacheDataMap[key] = CacheData{
		data:                     value,
		expirationEpochTimestamp: time.Now().Unix() + DefaultExpirationTimeInSeconds,
	}
}

// Get is used to get the value of a stored key in the cache. If the key is not exist, then it will return error (cache miss). If a key is exist, it will return the value
func (cc *CoreCache) Get(key string) (string, error) {
	cc.rwmu.RLock()
	defer cc.rwmu.RUnlock()

	cacheData, ok := cc.cacheDataMap[key]
	if !ok {
		return "", ErrorCacheMiss
	}

	return cacheData.data, nil

}

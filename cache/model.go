package cache

import (
	"errors"
	"sync"
)

// CacheData is the struct that represents the data that is stored in cacheDataMap for each key
type CacheData struct {
	data                     string
	expirationEpochTimestamp int64
}

// CoreCache is the struct that defines the core struct of the in-memory cache. It also holds the cacheDataMap in key-value schema.
type CoreCache struct {
	gracefulStopChannel chan struct{}

	wg           sync.WaitGroup
	rwmu         sync.RWMutex
	cacheDataMap map[string]CacheData
}

var (
	//ErrorCacheMiss this error happens when the key is not or no longer exist in cache
	ErrorCacheMiss error = errors.New("Key is not exist in Cache")

	//DefaultExpirationTimeInSeconds default expiration time for each key is in 30 minutes (1800 seconds)
	DefaultExpirationTimeInSeconds int64 = 1800
)

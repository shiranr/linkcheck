package cache_store

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

var cache *cacheStore

var fileName = "cache.json"

type CacheStore interface {
}

type cacheStore struct {
	CacheData []*CacheEntry
}

type CacheEntry struct {
	URL         string
	LastChecked time.Time
}

func GetCacheInstance() CacheStore {
	if cache == nil {
		cache = &cacheStore{}
		cache.loadCache()
	}
	return cache
}

func (cacheStore *cacheStore) loadCache() {
	file, _ := ioutil.ReadFile(fileName)
	err := json.Unmarshal([]byte(file), &cacheStore.CacheData)
	if err != nil {
		log.Errorf("Failed to load cache data falling to default")
	}
}

func (cacheStore *cacheStore) SaveToCache(url string) {
	cacheEntry := &CacheEntry{
		URL:         url,
		LastChecked: time.Now(),
	}
	cacheStore.CacheData = append(cacheStore.CacheData, cacheEntry)
	data, err := json.MarshalIndent(cache.CacheData, "", "")
	if err != nil {
		log.Errorf("Failed to marshal cache data with error %s", err)
		return
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Errorf("Failed to save cache file with error %s", err)
	}
}

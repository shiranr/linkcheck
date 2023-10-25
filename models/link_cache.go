package models

import (
	"encoding/gob"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

var cache *LinksCache
var filePath = "../cache_store/cache"

type LinksCache struct {
	linksCache map[string]*LinkData
	mapLock    sync.RWMutex
	duration   time.Duration
}

type LinkData struct {
	ResponseStatus int
	LastChecked    int64
	LinkPath       string
}

// Please notice this is not thread safe
func GetCacheInstance(customFilePath string, empty bool) *LinksCache {
	duration := viper.GetDuration("cache_duration")
	if cache == nil {
		if customFilePath != "" {
			filePath = customFilePath
		}
		cache = &LinksCache{
			linksCache: make(map[string]*LinkData),
			mapLock:    sync.RWMutex{},
			duration:   duration,
		}
		if !empty {
			cache.loadCacheData()
		}
	}
	return cache
}

func (c *LinksCache) Close() {
	c.SaveCache()
	cache = nil
}

func (c *LinksCache) loadCacheData() {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	c.decodeData(file)
}

func (c *LinksCache) SaveCache() {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	c.encodeData(file)
}

func (c *LinksCache) AddLink(linkPath string, status int) {
	c.mapLock.Lock()
	defer c.mapLock.Unlock()
	data := c.linksCache[linkPath]
	if data == nil {
		data = &LinkData{
			ResponseStatus: status,
			LastChecked:    time.Now().Unix(),
			LinkPath:       linkPath,
		}
	}
	c.linksCache[linkPath] = data
}

func (c *LinksCache) IsTimeCachedElapsed(linkPath string) bool {
	c.mapLock.RLock()
	defer c.mapLock.RUnlock()
	val, ok := c.linksCache[linkPath]
	if !ok {
		return true
	}
	return c.checkTimeElapsed(val)
}

func (c *LinksCache) checkTimeElapsed(val *LinkData) bool {
	if val.LastChecked+int64(c.duration.Seconds()) < time.Now().Unix() {
		return true
	}
	return false
}

func (c *LinksCache) CheckLinkStatus(linkPath string) (int, bool) {
	c.mapLock.RLock()
	defer c.mapLock.RUnlock()
	val, ok := c.linksCache[linkPath]
	if !ok {
		return 0, ok
	}
	if c.checkTimeElapsed(val) {
		return 0, false
	}
	return val.ResponseStatus, ok
}

func (c *LinksCache) encodeData(file *os.File) {
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(c.linksCache)
	if err != nil {
		panic(err)
	}
}

func (c *LinksCache) decodeData(file *os.File) {
	decoder := gob.NewDecoder(file)
	err := decoder.Decode(&c.linksCache)
	if err != nil {
		panic(err)
	}
}

package models

import (
	"encoding/gob"
	"os"
	"sync"
	"time"
)

var cache *LinksCache
var filePath = "../cache_store/cache"

type LinksCache struct {
	linksCache map[string]*LinkData
	mapLock    sync.RWMutex
}

type LinkData struct {
	ResponseStatus int
	LastChecked    int64
	LinkPath       string
}

// Please notice this is not thread safe
func GetCacheInstance(customFilePath string, empty bool) *LinksCache {
	if cache == nil {
		if customFilePath != "" {
			filePath = customFilePath
		}
		if empty {
			cache = &LinksCache{
				linksCache: make(map[string]*LinkData),
				mapLock:    sync.RWMutex{},
			}
		} else {
			cache = &LinksCache{}
			cache.loadCache()
		}
	}
	return cache
}

func (c *LinksCache) Close() {
	c.SaveCache()
	cache = nil
}

func (c *LinksCache) loadCache() {
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

func (c *LinksCache) CheckLinkCache(linkPath string) (int, bool) {
	c.mapLock.RLock()
	defer c.mapLock.RUnlock()
	val, ok := c.linksCache[linkPath]
	if !ok {
		return 0, ok
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

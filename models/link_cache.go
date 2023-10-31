package models

import (
	"encoding/gob"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var cache *LinksCache
var filePath string

type LinksCache struct {
	linksCache map[string]*LinkData
	mapLock    sync.RWMutex
	duration   time.Duration
	runNumber  uint64
}

type LinkData struct {
	ResponseStatus int
	LastChecked    int64
	LinkPath       string
	RunNumber      uint64
}

// Please notice this is not thread safe
func GetCacheInstance(empty bool) *LinksCache {
	if cache == nil {
		filePath = viper.GetString("cache_output_path")
		duration := viper.GetDuration("cache_duration")
		cache = &LinksCache{
			linksCache: make(map[string]*LinkData),
			mapLock:    sync.RWMutex{},
			duration:   duration,
			runNumber:  rand.Uint64(),
		}
		if !empty {
			log.Info("Loading cache data from " + filePath)
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
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		log.Error("File "+filePath+" does not exist, skipping cache load", err)
		return
	}
	log.Info("Opening cache file " + filePath)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	c.decodeData(file)
}

func (c *LinksCache) SaveCache() {
	basePath := filepath.Base(filePath)
	if err := os.MkdirAll(filepath.Dir(basePath), 0770); err != nil {
		log.Error("Failed to create path for cache file "+filePath, err)
		return
	}
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
			RunNumber:      c.runNumber,
		}
	} else {
		data.ResponseStatus = status
		data.RunNumber = c.runNumber
		data.LastChecked = time.Now().Unix()
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
	if val.RunNumber != c.runNumber && (299 < val.ResponseStatus || val.ResponseStatus < 200) {
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

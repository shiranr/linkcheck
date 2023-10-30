package tests

import (
	"github.com/shiranr/linkcheck/models"
	"github.com/shiranr/linkcheck/utils"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

var cache *models.LinksCache

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	configPath := basepath + "/resources/linkcheck.json"
	utils.LoadConfiguration(configPath)
	cache = models.GetCacheInstance(true)
	cache.SaveCache()
}

func TestInstancesAreTheSame(t *testing.T) {
	cache2 := models.GetCacheInstance(false)
	assert.Equal(t, cache, cache2)
}

func TestCacheIsNotNil(t *testing.T) {
	assert.NotNil(t, cache)
}

func TestAddingDataToCache(t *testing.T) {
	respStat, ok := cache.CheckLinkStatus("test")
	assert.Equal(t, respStat, 0)
	assert.False(t, ok)
	cache.AddLink("test", 200)
	cache.Close()
	cache = models.GetCacheInstance(false)
	respStat, ok = cache.CheckLinkStatus("test")
	assert.Equal(t, respStat, 200)
	assert.True(t, ok)
}

func TestCacheIsNotNilAfterClose(t *testing.T) {
	cache.Close()
	assert.NotNil(t, cache)
}

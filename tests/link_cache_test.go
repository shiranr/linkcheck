package tests

import (
	"github.com/shiranr/linkcheck/models"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

var cache *models.LinksCache

func getFilePath() string {
	basePath, _ := filepath.Abs(".")
	path := filepath.Join(basePath, "resources/test_cache")
	return path
}

func init() {
	path := getFilePath()
	cache = models.GetCacheInstance(path, true)
	cache.SaveCache()
}

func TestInstancesAreTheSame(t *testing.T) {
	path := getFilePath()
	cache2 := models.GetCacheInstance(path, false)
	assert.Equal(t, cache, cache2)
}

func TestCacheIsNotNil(t *testing.T) {
	assert.NotNil(t, cache)
}

func TestAddingDataToCache(t *testing.T) {
	respStat, ok := cache.CheckLinkCache("test")
	assert.Equal(t, respStat, 0)
	assert.False(t, ok)
	cache.AddLink("test", 200)
	cache.Close()
	cache = models.GetCacheInstance("resources/test_cache", false)
	respStat, ok = cache.CheckLinkCache("test")
	assert.Equal(t, respStat, 200)
	assert.True(t, ok)
}

func TestCacheIsNotNilAfterClose(t *testing.T) {
	cache.Close()
	assert.NotNil(t, cache)
}

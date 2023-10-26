package tests

import (
	"github.com/shiranr/linkcheck/models"
	"github.com/shiranr/linkcheck/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestFilesWithError(t *testing.T) {
	cache := models.GetCacheInstance("", true)
	validateLinkInCache(t, cache, "https://github.com/apache/jmeter", 0, false)
	readmeFiles := utils.ExtractMarkdownFiles()
	res := models.GetFilesProcessorInstance().Process(readmeFiles)
	assert.ErrorContains(t, res, "ERROR: 3 links check failed, please check the logs")
	logFilePath := setLogFile()
	logFileContent := readLogFile(logFilePath)
	assert.Contains(t, logFileContent, "Went through 2 files")
	assert.Contains(t, logFileContent, "Line 21 link nla.go status 400")
	assert.Contains(t, logFileContent, "Line 33 link source-control/merge-strategies.md status 400")
	assert.Contains(t, logFileContent, "Line 33 link resources/templates/CONTRIBUTING.md status 400")
	assert.NotContains(t, logFileContent, "http://bla.com/")
	assert.NotContains(t, logFileContent, "http://test.com/")
	cache = models.GetCacheInstance("", false)
	validateLinkInCache(t, cache, "http://bla.com/", 0, false)
	validateLinkInCache(t, cache, "http://test.com/", 0, false)
	validateLinkInCache(t, cache, "https://github.com/apache/jmeter", 200, true)

}

func validateLinkInCache(t *testing.T, cache *models.LinksCache, linkPath string, expectedStatus int, isOk bool) {
	status, ok := cache.CheckLinkStatus(linkPath)
	assert.Equal(t, ok, isOk)
	assert.Equal(t, status, expectedStatus)
}

func readLogFile(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func setLogFile() string {
	logFilePath := "resources/log.txt"
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).
			Fatal("Failed to open log file.")
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
	return logFilePath
}

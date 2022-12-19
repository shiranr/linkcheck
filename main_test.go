package main

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestLoadConfiguration(t *testing.T) {
	loadConfiguration()
	duration, _ := time.ParseDuration("30s")
	assert.Equal(t, viper.GetDuration("client_timeout"), duration)
}

func TestHandlerFile(t *testing.T) {
	loadConfiguration()
	files := extractReadmeFiles()
	assert.Len(t, files, 2)
	assert.Equal(t, files[0], "/shiranr/linkcheck/README.md")
	assert.Contains(t, files[1], "/shiranr/linkcheck/tests/resources/another_markdown.md")
}

func TestCheckLink(t *testing.T) {
	loadConfiguration()
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	handleFile(basepath + "/tests/resources/another_markdown.md")
	wg.Done()
	os.Open("")

}

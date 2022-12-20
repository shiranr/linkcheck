package main

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	loadConfiguration()
	assert.NotEmpty(t, viper.GetStringSlice("exclude_links"))
}

func TestHandlerFile(t *testing.T) {
	loadConfiguration()
	files := extractReadmeFiles()
	assert.Len(t, files, 2)
	assert.Contains(t, files[0], "/shiranr/linkcheck/README.md")
	assert.Contains(t, files[1], "/shiranr/linkcheck/tests/resources/another_markdown.md")
}

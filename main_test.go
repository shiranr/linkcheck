package main

import (
	"github.com/shiranr/linkcheck/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	utils.LoadConfiguration("configuration/linkcheck.json")
	assert.NotEmpty(t, viper.GetStringSlice("exclude_links"))
}

func TestHandlerFile(t *testing.T) {
	utils.LoadConfiguration("configuration/linkcheck.json")
	files := utils.ExtractMarkdownFiles()
	assert.Len(t, files, 2)
	assert.Contains(t, files[0], "/shiranr/linkcheck/README.md")
	assert.Contains(t, files[1], "/shiranr/linkcheck/tests/resources/another_markdown.md")
}

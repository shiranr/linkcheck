package utils

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigurationSuccessfully(t *testing.T) {
	LoadConfiguration("../configuration/linkcheck.json")
	assert.True(t, viper.GetBool("only_errors"))
}

func TestExtractMarkdownFilesSuccessfully(t *testing.T) {
	LoadConfiguration("../configuration/linkcheck.json")
	viper.Set("project_path", "../")
	files := ExtractMarkdownFiles()
	assert.Len(t, files, 3)
	assert.Contains(t, files[0], "/linkcheck/README.md")
	assert.Contains(t, files[1], "/linkcheck/tests/resources/MARKDOWN.md")
	assert.Contains(t, files[2], "/linkcheck/tests/resources/another_markdown.md")
}

func TestExtractMarkdownFilesFail(t *testing.T) {
	LoadConfiguration("../configuration/linkcheck.json")
	viper.Set("project_path", "")
	getDir = func() (dir string, err error) {
		return "", errors.New("BLA BLA BLA")
	}
	files := ExtractMarkdownFiles()
	assert.Len(t, files, 0)
}

func TestExtractMarkdownFilesFromList(t *testing.T) {
	LoadConfiguration("../configuration/linkcheck.json")
	files := []string{"../README.md", "test.jpg"}
	readmeFiles := ExtractMarkdownFilesFromList(files)
	assert.Len(t, readmeFiles, 1)
	assert.Contains(t, readmeFiles[0], "linkcheck/README.md")

}

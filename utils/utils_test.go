package utils

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	LoadConfiguration("../configuration/linkcheck.json")
	assert.True(t, viper.GetBool("only_errors"))
}

func TestHandlerFile(t *testing.T) {
	LoadConfiguration("../configuration/linkcheck.json")
	viper.Set("project_path", "../")
	files := ExtractMarkdownFiles()
	assert.Len(t, files, 3)
	assert.Contains(t, files[0], "/shiranr/linkcheck/README.md")
	assert.Contains(t, files[1], "/shiranr/linkcheck/tests/resources/MARKDOWN.md")
	assert.Contains(t, files[2], "/shiranr/linkcheck/tests/resources/another_markdown.md")
}

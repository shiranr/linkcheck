package models

import (
	"github.com/shiranr/linkcheck/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFilesProcessorInstance(t *testing.T) {
	fpi := GetFilesProcessorInstance()
	fpo := fpi.(*filesProcessor)
	assert.NotNil(t, fpo)
	assert.False(t, fpo.serial)
	assert.False(t, fpo.Result.done)
	assert.False(t, fpo.Result.close)
	assert.False(t, fpo.Result.onlyErrors)
	assert.NotNil(t, fpo.Result.FilesLinksMap)
	assert.NotNil(t, fpo.Result.mapLock)
	assert.NotNil(t, fpo.Result.Channel)
}

func TestProcessFiles(t *testing.T) {
	utils.LoadConfiguration("../configuration/linkcheck.json")
	files := []string{"../tests/resources/another_markdown.md"}
	fpi := GetFilesProcessorInstance()
	fpo := fpi.(*filesProcessor)
	fpo.linkProcessor = &mockLinkProcessor{
		t: t,
	}
	err := fpi.Process(files)
	assert.Nil(t, err)

}

type mockLinkProcessor struct {
	t *testing.T
}

func (mlp *mockLinkProcessor) CheckLink(filePath string, linkPath string, lineNumber int) *LinkResult {
	assert.Equal(mlp.t, filePath, "../tests/resources/another_markdown.md")
	assert.Equal(mlp.t, linkPath, "LINK")
	assert.Equal(mlp.t, lineNumber, 100000)
	return &LinkResult{
		lineNumber: 100000,
		status:     203,
		path:       "Hello/world",
		linkType:   "URL",
		filePath:   "PATH",
	}
}

func (mlp *mockLinkProcessor) ExtractLinks(fileData string) []*linkPath {
	assert.Contains(mlp.t, fileData, "[Test](MARKDOWN.md#energy-proportionality)")
	return []*linkPath{
		{
			LinkLineNumber: 100000,
			Link:           "LINK",
		},
	}
}

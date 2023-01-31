package models

import (
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

//func TestProcessFiles(t *testing.T) {
//	utils.LoadConfiguration("../configuration/linkcheck.json")
//	files := []string{"../README.md"}
//	fpi := GetFilesProcessorInstance()
//
//	err := fpi.Process(files)
//	assert.
//
//}

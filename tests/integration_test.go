package tests

import (
	"github.com/shiranr/linkcheck/models"
	"github.com/shiranr/linkcheck/utils"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFiles(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	configPath := basepath + "/resources/linkcheck.json"
	utils.LoadConfiguration(configPath)
	readmeFiles := utils.ExtractReadmeFiles()
	res := models.GetFilesProcessorInstance().Process(readmeFiles)
	assert.NoError(t, res)
}

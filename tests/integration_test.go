package tests

import (
	"github.com/shiranr/linkcheck/models"
	"github.com/shiranr/linkcheck/utils"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFilesWithError(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	configPath := basepath + "/resources/linkcheck.json"
	utils.LoadConfiguration(configPath)
	readmeFiles := utils.ExtractReadmeFiles()
	res := models.GetFilesProcessorInstance().Process(readmeFiles)
	assert.ErrorContains(t, res, "ERROR: 3 links check failed, please check the logs")

}

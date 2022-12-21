package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func LoadConfiguration(configPath string) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to load configuration")
	}
}

func ExtractReadmeFiles() []string {
	var readmeFiles []string
	path, err := os.Getwd()
	if err == nil {
		err = filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
			if file.IsDir() && strings.Contains(file.Name(), "vendor") {
				return filepath.SkipDir
			}
			if strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
				path, _ = filepath.Abs(path)
				readmeFiles = append(readmeFiles, path)
			}
			return nil
		})
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get files")
	}
	return readmeFiles
}

func ExtractReadmeFilesFromList(filesList []string) []string {
	var readmeFiles []string
	for _, filePath := range filesList {
		if strings.HasSuffix(strings.ToLower(filePath), ".md") {
			filePath, _ = filepath.Abs(filePath)
			readmeFiles = append(readmeFiles, filePath)
		}
	}
	return readmeFiles
}

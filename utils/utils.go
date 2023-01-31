package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var getDir = os.Getwd

// LoadConfiguration - load configuration file from config path
func LoadConfiguration(configPath string) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to load configuration")
	}
}

// ExtractMarkdownFiles - extract markdown files from the defined project_path or current dir we are on
//and walk through the path
func ExtractMarkdownFiles() []string {
	var markdownFiles []string
	var err error
	path := viper.GetString("project_path")
	if path == "" {
		path, err = getDir()
	}
	if err == nil {
		log.Info("extracting markdown files from path " + path)
		err = filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("Failed to walk over path: " + path)
				return err
			}
			if file.IsDir() && strings.Contains(file.Name(), "vendor") {
				return filepath.SkipDir
			}
			if strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
				path, _ = filepath.Abs(path)
				markdownFiles = append(markdownFiles, path)
			}
			return nil
		})
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to get files")
	}
	return markdownFiles
}

// ExtractMarkdownFilesFromList - extract markdown files from given list
func ExtractMarkdownFilesFromList(filesList []string) []string {
	var markdownFiles []string
	log.Info("extracting markdown files from list")
	for _, filePath := range filesList {
		if strings.HasSuffix(strings.ToLower(filePath), ".md") {
			filePath, _ = filepath.Abs(filePath)
			markdownFiles = append(markdownFiles, filePath)
		}
	}
	return markdownFiles
}

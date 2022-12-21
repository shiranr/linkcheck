package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func SetUpLogger() {
	outputPath := viper.GetString("output_path")
	logFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).
			Fatal("Failed to open log file.")
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
}

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
	path := viper.GetString("path")
	var readmeFiles []string
	if envPath := os.Getenv("PROJECT_PATH"); envPath != "" {
		path = envPath
	}
	err := filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() && strings.Contains(file.Name(), "vendor") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
			path, _ = filepath.Abs(path)
			readmeFiles = append(readmeFiles, path)
		}
		return nil
	})

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

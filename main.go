package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"linkcheck/models"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// TODO add tests.
// TODO add CMD.
// TODO make this a linter for megalinter.
// TODO add workflow
func main() {
	start := time.Now()
	loadConfiguration()
	outputPath := viper.GetString("output_path")
	logFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).
			Fatal("Failed to open log file.")
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
	readmeFiles := extractReadmeFiles()
	models.GetFilesProcessorInstance().Process(readmeFiles)
	end := time.Now()
	log.Info("Time elapsed: " + end.Sub(start).String())
}

func loadConfiguration() {
	viper.SetConfigName("linkcheck.json")
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	viper.AddConfigPath(basepath + "/configuration")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to load configuration")
	}
}

func extractReadmeFiles() []string {
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

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
	"sync"
	"time"
)

var wg sync.WaitGroup
var lineHandler = models.GetLinkExtractorInstance()

var result = models.Result{
	FilesLinksMap: map[string]*models.FileLink{},
}

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
	extractLinksFromReadmes(readmeFiles)
	wg.Wait()
	result.Print()
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

func extractLinksFromReadmes(files []string) {
	for _, filePath := range files {
		wg.Add(1)
		go handleFile(filePath)
	}
}

func handleFile(filePath string) {
	defer wg.Done()
	fileLinkData := models.FileLink{
		FilePath: filePath,
		Links:    []*models.Link{},
	}
	result.AddNewFile(&fileLinkData)
	fileData, err := models.NewFileData(filePath)
	if err != nil {
		return
	}
	lineText, lineNumber := fileData.ScanOneLine()
	for lineNumber != -1 {
		linksPaths := lineHandler.FindAndGetLinks(lineText)
		for _, linkPath := range linksPaths {
			linkData := &models.Link{
				LineNumber: lineNumber,
				Status:     0,
				Path:       linkPath,
			}
			wg.Add(1)
			go checkLink(fileData, linkData)
		}
		lineText, lineNumber = fileData.ScanOneLine()
	}
}

func checkLink(fileData *models.FileData, linkData *models.Link) {
	defer wg.Done()
	switch {
	case strings.HasPrefix(linkData.Path, "http"):
		urlHandler := models.GetURLHandlerInstance()
		urlHandler.Handle(linkData)
	case strings.HasPrefix(linkData.Path, "mailto:"):
		emailHandler := models.GetEmailHandlerInstance()
		emailHandler.Handle(linkData)
	default:
		fileLinkHandler := models.GetFileLinkHandler(fileData)
		fileLinkHandler.Handle(linkData)
	}
	result.Append(linkData, fileData.FullFilePath())
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

package models

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// FileProcessor  - process a single file, get links and check them
type FileProcessor interface {
	ProcessFile()
}

type fileProcessor struct {
	resultChan    chan *LinkResult
	filePath      string
	fileLinesData map[int]string
}

// GetNewFileProcessor - get instance of file processor
func GetNewFileProcessor(filePath string, resultChan chan *LinkResult) FileProcessor {
	return &fileProcessor{
		filePath:   filePath,
		resultChan: resultChan,
	}
}

// ProcessFile - process a single file
func (fp *fileProcessor) ProcessFile() {
	defer wg.Done()
	linkHandler := GetLinkProcessorInstance()
	fileBytes, err := os.ReadFile(fp.filePath)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to read file " + fp.filePath)
		linkData := &LinkResult{
			lineNumber: 0,
			status:     404,
			filePath:   fp.filePath,
			linkType:   InternalLink,
		}
		fp.resultChan <- linkData
		return
	}
	fileData := string(fileBytes)
	linksPaths := linkHandler.ExtractLinks(fileData)
	for _, linkData := range linksPaths {
		fp.resultChan <- linkHandler.CheckLink(fp.filePath, linkData.Link, linkData.LinkLineNumber)
	}
}

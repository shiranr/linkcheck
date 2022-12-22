package models

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type FileHandler interface {
	HandleFile()
}

type fileHandler struct {
	resultChan    chan *LinkResult
	filePath      string
	fileLinesData map[int]string
}

func GetNewFileHandler(filePath string, resultChan chan *LinkResult) FileHandler {
	return &fileHandler{
		filePath:   filePath,
		resultChan: resultChan,
	}
}

func (fh *fileHandler) HandleFile() {
	defer wg.Done()
	linkHandler := GetLinkHandlerInstance()
	fileBytes, err := os.ReadFile(fh.filePath)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to read file " + fh.filePath)
		linkData := &Link{
			LineNumber: 0,
			Status:     404,
			Path:       fh.filePath,
			LinkType:   File,
		}
		fh.resultChan <- &LinkResult{filePath: fh.filePath, link: linkData}
		return
	}
	fileData := string(fileBytes)
	linksPaths := linkHandler.ExtractLinks(fileData)
	for _, linkData := range linksPaths {
		linkData := linkHandler.CheckLink(fh.filePath, linkData.Link, linkData.LinkLineNumber)
		fh.resultChan <- &LinkResult{filePath: fh.filePath, link: linkData}
	}
}

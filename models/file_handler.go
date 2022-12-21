package models

import (
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
	fileBytes, _ := os.ReadFile(fh.filePath)
	fileData := string(fileBytes)
	linksPaths := linkHandler.ExtractLinks(fileData)
	for _, linkData := range linksPaths {
		linkData := linkHandler.CheckLink(fh.filePath, linkData.Link, linkData.LinkLineNumber)
		fh.resultChan <- &LinkResult{filePath: fh.filePath, link: linkData}
	}
}

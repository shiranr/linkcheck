package models

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
)

type FileHandler interface {
	HandleFile()
}

type fileHandler struct {
	resultChan chan *LinkResult
	filePath   string
	Scanner
}

type Scanner struct {
	scanner    *bufio.Scanner
	lineNumber int
}

func GetNewFileHandler(filePath string, resultChan chan *LinkResult) FileHandler {
	file, err := os.Open(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"path":  filePath,
		}).Errorf("Failed to read file.")
		return nil
	}
	return &fileHandler{
		filePath: filePath,
		Scanner: Scanner{
			scanner:    bufio.NewScanner(file),
			lineNumber: 1,
		},
		resultChan: resultChan,
	}
}

func (fh *fileHandler) HandleFile() {
	defer wg.Done()
	linkHandler := GetLinkHandlerInstance()
	lineText, lineNumber := fh.scanOneLine()

	for lineNumber != -1 {
		linksPaths := linkHandler.ExtractLinks(lineText)
		for _, linkPath := range linksPaths {
			linkData := linkHandler.CheckLink(fh.filePath, linkPath, lineNumber)
			fh.resultChan <- &LinkResult{filePath: fh.filePath, link: linkData}
		}
		lineText, lineNumber = fh.scanOneLine()
	}
}

func (fh *fileHandler) scanOneLine() (string, int) {
	if fh.scanner.Scan() {
		lineText := fh.scanner.Text()
		fh.lineNumber++
		return lineText, fh.lineNumber
	}
	return "", -1
}

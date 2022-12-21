package models

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	linkFileHandler *fileLinkHandler
)

type fileLinkHandler struct {
	folderPath string
	fileName   string
}

func GetFileLinkHandler(filePath string) *fileLinkHandler {
	fileData, err := NewFileData(filePath)
	if err != nil {
		return nil
	}
	linkFileHandler = &fileLinkHandler{
		folderPath: fileData.folderPath,
		fileName:   fileData.fileName,
	}
	return linkFileHandler
}

func (handler *fileLinkHandler) Handle(linkPath string) int {
	linkedFileEscapedFullPath := handler.escapedFullPath(linkPath)
	_, err := os.Stat(linkedFileEscapedFullPath)
	if err != nil {
		log.WithFields(log.Fields{
			"link":  linkPath,
			"error": err,
		}).Error("Failed to get link data")
		return 400
	}
	if strings.Contains(linkPath, "#") {
		fileBytes, _ := os.ReadFile(linkedFileEscapedFullPath)
		fileData := string(fileBytes)
		if !handler.fileContainsLink(linkPath, fileData) {
			return 400
		}
	}
	return 200
}

func (handler *fileLinkHandler) escapedFullPath(extension string) string {
	folderPath := handler.folderPath
	if strings.HasPrefix(extension, "#") {
		folderPath = path.Join(handler.folderPath, handler.fileName)
	} else if strings.Contains(extension, "#") {
		fileName := strings.Split(extension, "#")[0]
		folderPath = path.Join(handler.folderPath, fileName)
	} else {
		folderPath = path.Join(handler.folderPath, extension)
	}
	folderPath, _ = url.PathUnescape(folderPath)
	return folderPath
}

func (handler *fileLinkHandler) fileContainsLink(titleLink string, fileText string) bool {
	titleLink = strings.Split(titleLink, "#")[1]
	title := strings.ReplaceAll(titleLink, "#", "")
	title = strings.ReplaceAll(title, "-", "( |-|)")
	readmeTitleRegex := "(?i)#( ?)" + title
	linkRegex, _ := regexp.Compile(readmeTitleRegex)
	if len(linkRegex.FindStringSubmatch(fileText)) > 0 {
		return true
	}
	return false
}

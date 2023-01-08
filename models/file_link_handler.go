package models

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	internalLink *internalLinkHandler
)

type internalLinkHandler struct {
	filePath string
}

func GetInternalLinkHandler(filePath string) LinkHandlerInterface {
	internalLink = &internalLinkHandler{
		filePath: filePath,
	}
	return internalLink
}

func (handler *internalLinkHandler) Handle(linkPath string) int {
	folderPath, fileName := filepath.Split(handler.filePath)
	linkedFileEscapedFullPath := handler.escapedFullPath(folderPath, fileName, linkPath)
	if strings.Contains(linkedFileEscapedFullPath, "?") {
		linkedFileEscapedFullPath = strings.Split(linkedFileEscapedFullPath, "?")[0]
	}
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

func (handler *internalLinkHandler) escapedFullPath(folderPath, fileName, linkPath string) string {
	if strings.HasPrefix(linkPath, "#") {
		folderPath = path.Join(folderPath, fileName)
	} else if strings.Contains(linkPath, "#") {
		fileName := strings.Split(linkPath, "#")[0]
		folderPath = path.Join(folderPath, fileName)
	} else {
		folderPath = path.Join(folderPath, linkPath)
	}
	folderPath, _ = url.PathUnescape(folderPath)
	return folderPath
}

func (handler *internalLinkHandler) fileContainsLink(titleLink string, fileText string) bool {
	titleLink = strings.Split(titleLink, "#")[1]
	title := strings.ReplaceAll(titleLink, "#", "")
	title = strings.ReplaceAll(title, "(", "\\(")
	title = strings.ReplaceAll(title, ")", "\\)")
	title = "(?i)" + title
	title = strings.ReplaceAll(title, "-", "( |-||: )(?i)")
	readmeTitleRegex := "#(.*)(?i)" + title
	linkRegex, err := regexp.Compile(readmeTitleRegex)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "title_link": titleLink}).Error("Failed to create file link regex")
		return false
	}
	if len(linkRegex.FindStringSubmatch(fileText)) > 0 {
		return true
	}
	return false
}

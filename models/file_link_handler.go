package models

import (
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	fileHandler *fileLinkHandler
)

type fileLinkHandler struct {
	fileData *FileData
}

func GetfileLinkHandler(data *FileData) *fileLinkHandler {
	fileHandler = &fileLinkHandler{
		fileData: data,
	}
	return fileHandler
}

func (handler *fileLinkHandler) Handle(linkData *Link) {
	linkData.LinkType = Folder
	linkedFileEscapedFullPath := handler.escapedFullPath(linkData.Path)
	_, err := os.Stat(linkedFileEscapedFullPath)
	if err != nil {
		linkData.Status = 400
		println("Failed to get link data with path " + linkData.Path + " and error " + err.Error())
		return
	}
	if strings.Contains(linkData.Path, "#") {
		fileBytes, _ := os.ReadFile(linkedFileEscapedFullPath)
		fileData := string(fileBytes)
		if !handler.fileContainsLink(linkData.Path, fileData) {
			linkData.Status = 400
			return
		}
	}
	linkData.Status = 200
}

func (handler *fileLinkHandler) escapedFullPath(extension string) string {
	folderPath := handler.fileData.folderPath
	if strings.HasPrefix(extension, "#") {
		folderPath = path.Join(handler.fileData.folderPath, handler.fileData.fileName)
	} else if strings.Contains(extension, "#") {
		fileName := strings.Split(extension, "#")[0]
		folderPath = path.Join(handler.fileData.folderPath, fileName)
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

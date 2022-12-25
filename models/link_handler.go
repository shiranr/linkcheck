package models

import (
	"github.com/spf13/viper"
	"regexp"
	"strings"
	"sync"
)

var link = "\\[{1}([é’*&\"|`?'>\\-\\sa-zA-Z0-9@:%._\\\\+~#=,\\n\\/\\(\\)])*((\\]\\()){1}([?\\sa-zA-Z0-9@:%._\\\\+~#=\\/\\/\\-]{1,256}(\\(.*\\))?(\\\"(.*)\\\")?)\\){1}"

var urlRegex = "https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"

type LinkHandler interface {
	CheckLink(filePath string, linkPath string, lineNumber int) *Link
	ExtractLinks(fileData string) []*linkPath
}

type linkHandler struct {
	linksCache
	mdLinkRegex   *regexp.Regexp
	excludedLinks []string
}

type linksCache struct {
	linksCache map[string]int
	mapLock    sync.RWMutex
}

func (cache *linksCache) addLink(linkPath string, status int) {
	cache.mapLock.Lock()
	defer cache.mapLock.Unlock()
	cache.linksCache[linkPath] = status
}

func (cache *linksCache) readLink(linkPath string) (int, bool) {
	cache.mapLock.RLock()
	defer cache.mapLock.RUnlock()
	val, ok := cache.linksCache[linkPath]
	return val, ok
}

var lh *linkHandler

func GetLinkHandlerInstance() LinkHandler {
	linkOrPath := link + "|" + urlRegex
	regex, _ := regexp.Compile(linkOrPath)
	if lh == nil {
		lh = &linkHandler{
			linksCache{
				linksCache: map[string]int{},
			},
			regex,
			viper.GetStringSlice("exclude_links"),
		}
	}
	return lh
}

func (handler *linkHandler) CheckLink(filePath string, linkPath string, lineNumber int) *Link {
	linkData := &Link{
		LineNumber: lineNumber,
		Status:     0,
		Path:       linkPath,
	}
	status, ok := handler.readLink(linkPath)
	if !ok {
		switch {
		case strings.HasPrefix(linkData.Path, "http"):
			linkData.LinkType = URL
			urlHandler := GetURLHandlerInstance()
			linkData.Status = urlHandler.Handle(linkPath)
		case strings.HasPrefix(linkData.Path, "mailto:"):
			linkData.LinkType = Email
			emailHandler := GetEmailHandlerInstance()
			linkData.Status = emailHandler.Handle(linkPath)
		default:
			linkData.LinkType = Folder
			fileLinkHandler := GetFileLinkHandler(filePath)
			linkData.Status = fileLinkHandler.Handle(linkPath)
		}
		handler.addLink(linkPath, linkData.Status)
	} else {
		linkData.Status = status
	}
	return linkData
}

type linkPath struct {
	LinkLineNumber int
	Link           string
}

func (handler *linkHandler) ExtractLinks(fileData string) []*linkPath {
	var readmeLinks []string
	var validLinks []*linkPath
	readmeLinks = append(readmeLinks, handler.mdLinkRegex.FindAllString(fileData, -1)...)
	for _, path := range readmeLinks {
		if strings.Contains(path, "](") {
			path = strings.Split(path, "](")[1]
			if strings.HasSuffix(path, "))") {
				path = path[0 : len(path)-1]
			}
			if strings.HasSuffix(path, ")") && !strings.Contains(path, "(") {
				path = path[0 : len(path)-1]
			}
			path = strings.Split(path, " \\'")[0]
			path = strings.Split(path, " \"")[0]
		}
		if !handler.isExcluded(path) {
			linkPath := &linkPath{LinkLineNumber: handler.findLineNumber(path, fileData), Link: path}
			validLinks = append(validLinks, linkPath)
		}
	}
	return validLinks
}

func (handler *linkHandler) findLineNumber(link string, fileData string) int {
	for index, line := range strings.Split(fileData, "\n") {
		if strings.Contains(line, link) {
			return index + 1
		}
	}
	return -1
}

func (handler *linkHandler) isExcluded(link string) bool {
	for _, excludedPath := range viper.GetStringSlice("exclude_links") {
		if strings.HasPrefix(link, excludedPath) {
			return true
		}
	}
	return false
}

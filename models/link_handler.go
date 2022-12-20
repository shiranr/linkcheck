package models

import (
	"github.com/spf13/viper"
	"regexp"
	"strings"
	"sync"
)

var link = "\\[.*\\]\\([-a-zA-Z0-9@:%_\\+.~#?&=\\s]\\)|https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"

type LinkHandler interface {
	CheckLink(filePath string, linkPath string, lineNumber int) *Link
	ExtractLinks(path string) []string
}

type linkHandler struct {
	linksCache
	regex         *regexp.Regexp
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
	if lh == nil {
		var regex, _ = regexp.Compile(link)
		lh = &linkHandler{
			linksCache{
				linksCache: map[string]int{},
			}, regex, viper.GetStringSlice("exclude_links"),
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

func (handler *linkHandler) ExtractLinks(path string) []string {
	var linksPaths []string
	var validPaths []string
	linksPaths = append(linksPaths, handler.regex.FindAllString(path, -1)...)
	for _, linkPath := range linksPaths {
		if strings.Contains(linkPath, "](") {
			linkPath = strings.Split(linkPath, "](")[1]
		}
		lastIndex := strings.LastIndex(linkPath, ")")
		if lastIndex > 0 {
			linkPath = linkPath[0:lastIndex]
		}
		if !handler.isExcluded(linkPath) {
			validPaths = append(validPaths, linkPath)
		}
	}
	return validPaths
}

func (handler *linkHandler) isExcluded(link string) bool {
	for _, excludedPath := range viper.GetStringSlice("exclude_links") {
		if strings.HasPrefix(link, excludedPath) {
			return true
		}
	}
	return false
}

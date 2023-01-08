package models

import (
	"github.com/spf13/viper"
	"regexp"
	"strings"
	"sync"
)

var link = "\\[{1}([é’*&\"|`?'>\\-\\sa-zA-Z0-9@:%._\\\\+~#=,\\n\\/\\(\\)])*((\\]\\()){1}([?\\sa-zA-Z0-9@:%._\\\\+~#=\\/\\/\\-]{1,256}(\\(.*\\))?(\\\"(.*)\\\")?)\\){1}"

var urlRegex = "https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"

// LinkProcessor - process a single link, sends it to handler according to its type
type LinkProcessor interface {
	CheckLink(filePath string, linkPath string, lineNumber int) *LinkResult
	ExtractLinks(fileData string) []*linkPath
}

type linkProcessor struct {
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

func (cache *linksCache) checkLinkCache(linkPath string) (int, bool) {
	cache.mapLock.RLock()
	defer cache.mapLock.RUnlock()
	val, ok := cache.linksCache[linkPath]
	return val, ok
}

var lh *linkProcessor

// GetLinkProcessorInstance - get instance of link processor
func GetLinkProcessorInstance() LinkProcessor {
	linkOrPath := link + "|" + urlRegex
	regex, _ := regexp.Compile(linkOrPath)
	if lh == nil {
		lh = &linkProcessor{
			linksCache{
				linksCache: map[string]int{},
			},
			regex,
			viper.GetStringSlice("exclude_links"),
		}
	}
	return lh
}

// CheckLink - check a single link and pass it to its handler according to the type of the link
func (processor *linkProcessor) CheckLink(filePath string, linkPath string, lineNumber int) *LinkResult {
	linkData := &LinkResult{
		lineNumber: lineNumber,
		status:     0,
		path:       linkPath,
		filePath:   filePath,
	}
	status, ok := processor.checkLinkCache(linkPath)
	if !ok {
		switch {
		case strings.HasPrefix(linkData.path, "http"):
			linkData.linkType = URL
			urlHandler := GetURLHandlerInstance()
			linkData.status = urlHandler.Handle(linkPath)
		case strings.HasPrefix(linkData.path, "mailto:"):
			linkData.linkType = Email
			emailHandler := GetEmailHandlerInstance()
			linkData.status = emailHandler.Handle(linkPath)
		default:
			linkData.linkType = InternalLink
			fileLinkHandler := GetInternalLinkHandler(filePath)
			linkData.status = fileLinkHandler.Handle(linkPath)
		}
		processor.addLink(linkPath, linkData.status)
	} else {
		linkData.status = status
	}
	return linkData
}

type linkPath struct {
	LinkLineNumber int
	Link           string
}

// ExtractLinks - extract all the links from a single file
func (processor *linkProcessor) ExtractLinks(fileData string) []*linkPath {
	var readmeLinks []string
	var validLinks []*linkPath
	readmeLinks = append(readmeLinks, processor.mdLinkRegex.FindAllString(fileData, -1)...)
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
		if !processor.isExcluded(path) {
			linkPath := &linkPath{LinkLineNumber: processor.findLineNumber(path, fileData), Link: path}
			validLinks = append(validLinks, linkPath)
		}
	}
	return validLinks
}

func (processor *linkProcessor) findLineNumber(link string, fileData string) int {
	for index, line := range strings.Split(fileData, "\n") {
		if strings.Contains(line, link) {
			return index + 1
		}
	}
	return -1
}

func (processor *linkProcessor) isExcluded(link string) bool {
	for _, excludedPath := range viper.GetStringSlice("exclude_links") {
		if strings.HasPrefix(link, excludedPath) {
			return true
		}
	}
	return false
}

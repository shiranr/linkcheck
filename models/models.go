package models

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type FileLink struct {
	FilePath string
	Links    []*Link
}

func (fileLink *FileLink) append(link *Link) {
	fileLink.Links = append(fileLink.Links, link)
}

type LinkType string

const (
	URL    LinkType = "URL"
	Email           = "Email"
	Folder          = "Folder"
)

type Link struct {
	LineNumber int
	Status     int
	Path       string
	LinkType   LinkType
}

type Result struct {
	FilesLinksMap map[string]*FileLink
	mapLock       sync.RWMutex
	Channel       chan *LinkResult
	close         bool
	done          bool
}

type LinkResult struct {
	filePath string
	link     *Link
}

func getResult() *Result {
	result := &Result{
		FilesLinksMap: map[string]*FileLink{},
		Channel:       make(chan *LinkResult),
		close:         false,
		done:          false,
	}
	go result.Read()
	return result
}

func (res *Result) AddNewFile(fileLink *FileLink) {
	res.mapLock.Lock()
	defer res.mapLock.Unlock()
	res.FilesLinksMap[fileLink.FilePath] = fileLink
}

func (res *Result) Read() {
	for !res.close || len(res.Channel) > 0 {
		linkResult := <-res.Channel
		if linkResult != nil {
			res.Append(linkResult.link, linkResult.filePath)
		}
	}
	res.done = true
}

func (res *Result) Close() {
	res.close = true
	close(res.Channel)
	for !res.done {
	}
}

func (res *Result) Append(link *Link, filePath string) {
	res.mapLock.Lock()
	defer res.mapLock.Unlock()
	res.FilesLinksMap[filePath].append(link)
}

func (res *Result) Print() {
	log.Info("Went through " + strconv.Itoa(len(res.FilesLinksMap)) + " links")
	for key, val := range res.FilesLinksMap {
		log.Info("****************************")
		log.Info("Results for file " + key)
		log.Info("")
		for _, link := range val.Links {
			log.Info("Line " + strconv.Itoa(link.LineNumber) + " link " + link.Path + " status " + strconv.Itoa(link.Status))
			log.Info("")
		}
	}
}

package models

import (
	"strconv"
	"sync"
)

type FileLink struct {
	FilePath string
	Links    []Link
}

func (fileLink *FileLink) append(link Link) {
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
	MapLock       sync.RWMutex
}

func (result *Result) AddNewFile(fileLink *FileLink) {
	result.MapLock.Lock()
	defer result.MapLock.Unlock()
	result.FilesLinksMap[fileLink.FilePath] = fileLink
}

func (result *Result) Append(link Link, filePath string) {
	result.MapLock.Lock()
	defer result.MapLock.Unlock()
	result.FilesLinksMap[filePath].append(link)
}

func (result *Result) Print() {
	for key, val := range result.FilesLinksMap {
		println("****************************")
		println("Results for file " + key)
		println("")
		for _, link := range val.Links {
			println("Line " + strconv.Itoa(link.LineNumber) + " link " + link.Path + " status " + strconv.Itoa(link.Status))
			println("")
		}
	}
}

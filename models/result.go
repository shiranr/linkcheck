package models

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"sync"
)

// FileResultData - contains all file links and their data
type FileResultData struct {
	FilePath string
	Error    bool
	Links    []*LinkResult
}

func (frd *FileResultData) append(link *LinkResult) {
	frd.Links = append(frd.Links, link)
}

// LinkType - 4 types of links we currently have url/email/file
type LinkType string

// const - Enum to determine the type of link
const (
	URL          LinkType = "URL"
	Email                 = "Email"
	InternalLink          = "InternalLink"
)

// LinkResult - a descriptor which contains data for a single link
type LinkResult struct {
	lineNumber int
	status     int
	path       string
	linkType   LinkType
	filePath   string
}

// Result - a result for an entire run use channel to get results
type Result struct {
	FilesLinksMap map[string]*FileResultData
	mapLock       sync.RWMutex
	Channel       chan *LinkResult
	close         bool
	done          bool
	onlyErrors    bool
}

func getResult() *Result {
	result := &Result{
		FilesLinksMap: map[string]*FileResultData{},
		Channel:       make(chan *LinkResult),
		close:         false,
		done:          false,
		onlyErrors:    viper.GetBool("only_errors"),
	}
	go result.read()
	return result
}

// AddNewFile - add new file to the FileLinksMap to be saved as part of the results
func (res *Result) AddNewFile(fileLink *FileResultData) {
	res.mapLock.Lock()
	defer res.mapLock.Unlock()
	res.FilesLinksMap[fileLink.FilePath] = fileLink
}

func (res *Result) read() {
	for !res.close || len(res.Channel) > 0 {
		linkResult := <-res.Channel
		if linkResult != nil {
			res.append(linkResult, linkResult.filePath)
		}
	}
	res.done = true
}

func (res *Result) append(link *LinkResult, filePath string) {
	res.mapLock.Lock()
	defer res.mapLock.Unlock()
	fileData := res.FilesLinksMap[filePath]
	if link.status != 200 {
		fileData.Error = true
	}
	res.FilesLinksMap[filePath].append(link)
}

// Close = close the channel when we are done with all results.
func (res *Result) Close() {
	res.close = true
	close(res.Channel)
	for !res.done {
	}
}

// Print - print and return the results and return the errors we received
func (res *Result) Print() error {
	errCount := 0
	log.Info("Went through " + strconv.Itoa(len(res.FilesLinksMap)) + " files")
	for key, val := range res.FilesLinksMap {
		if !res.onlyErrors || res.onlyErrors && val.Error {
			log.Info("****************************")
			log.Info("Results for file " + key)
			log.Info("")
			for _, link := range val.Links {
				if !res.onlyErrors || res.onlyErrors && link.status != 200 {
					if link.status != 200 {
						errCount++
					}
					log.Info("Line " + strconv.Itoa(link.lineNumber) + " link " + link.path + " status " + strconv.Itoa(link.status))
					log.Info("")
				}
			}
		}
	}
	if errCount > 0 {
		errMsg := "ERROR: " + strconv.Itoa(errCount) + " links check failed, please check the logs"
		return errors.New(errMsg)
	}
	return nil
}

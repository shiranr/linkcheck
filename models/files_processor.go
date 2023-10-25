package models

import (
	"github.com/spf13/viper"
	"sync"
)

var wg sync.WaitGroup

var fp *filesProcessor

type filesProcessor struct {
	*Result
	serial bool
}

// FilesProcessor - process multiple files either parallel or one by one
type FilesProcessor interface {
	Process(files []string) error
}

// GetFilesProcessorInstance - get instance of files processor (Singleton)
func GetFilesProcessorInstance() FilesProcessor {
	if fp == nil {
		fp = &filesProcessor{
			getResult(),
			viper.GetBool("serial"),
		}
	}
	return fp
}

// Process - process the multiple files list
func (fh *filesProcessor) Process(files []string) error {
	for _, filePath := range files {
		fileLinkData := FileResultData{
			FilePath: filePath,
			Links:    []*LinkResult{},
			Error:    false,
		}
		fh.AddNewFile(&fileLinkData)
		fileProcessor := GetNewFileProcessor(filePath, fh.Channel)
		if fileProcessor != nil {
			wg.Add(1)
			fh.invoke(fileProcessor)
		}
	}
	wg.Wait()
	fh.Close()
	cache := GetCacheInstance("", false)
	cache.Close()
	return fh.Print()
}

func (fh *filesProcessor) invoke(fileProcessor FileProcessor) {
	if fh.serial {
		fileProcessor.ProcessFile()
	} else {
		go fileProcessor.ProcessFile()
	}
}

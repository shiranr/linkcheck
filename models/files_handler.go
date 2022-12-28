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

type FilesProcessor interface {
	Process(files []string) error
}

func GetFilesProcessorInstance() FilesProcessor {
	if fp == nil {
		fp = &filesProcessor{
			getResult(),
			viper.GetBool("serial"),
		}
	}
	return fp
}

func (fh *filesProcessor) Process(files []string) error {
	for _, filePath := range files {
		fileLinkData := FileLink{
			FilePath: filePath,
			Links:    []*Link{},
			Error:    false,
		}
		fh.AddNewFile(&fileLinkData)
		fileHandler := GetNewFileHandler(filePath, fh.Channel)
		if fileHandler != nil {
			wg.Add(1)
			fh.invoke(fileHandler)
		}
	}
	wg.Wait()
	fh.Close()
	return fh.Print()
}

func (fh *filesProcessor) invoke(fileHandler FileHandler) {
	if fh.serial {
		fileHandler.HandleFile()
	} else {
		go fileHandler.HandleFile()
	}
}

package models

import "sync"

var wg sync.WaitGroup

var fp *filesProcessor

type filesProcessor struct {
	*Result
}

type FilesProcessor interface {
	Process(files []string) error
}

func GetFilesProcessorInstance() FilesProcessor {
	if fp == nil {
		fp = &filesProcessor{
			getResult(),
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
			go fileHandler.HandleFile()
		}
	}
	wg.Wait()
	fh.Close()
	return fh.Print()
}

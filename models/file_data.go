package models

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
)

type FileData struct {
	fileName   string
	folderPath string
	file       *os.File
}

func (fileData *FileData) FileName() string {
	return fileData.fileName
}

func (fileData *FileData) FolderPath() string {
	return fileData.folderPath
}

func (fileData *FileData) FullFilePath() string {
	return path.Join(fileData.folderPath, fileData.fileName)
}

func (fileData *FileData) File() *os.File {
	return fileData.file
}

func NewFileData(path string) (*FileData, error) {
	fileData := &FileData{}
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.WithFields(log.Fields{
			"path":  path,
			"error": err,
		}).Error("Failed to read file")
		return nil, err
	}
	if fileInfo.IsDir() {
		fileData.folderPath = path
	} else {
		fileData.folderPath, _ = filepath.Split(path)
		fileData.fileName = filepath.Base(path)
	}
	return fileData, nil

}

func (fileData *FileData) Close() {
	defer fileData.file.Close()
}

package models

import (
	"bufio"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileData struct {
	fileName   string
	folderPath string
	file       *os.File
	Scanner
}

type Scanner struct {
	scanner    *bufio.Scanner
	lineNumber int
	canRead    bool
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

func (fileData *FileData) EscapedFullPath(extension string) string {
	folderPath := fileData.folderPath
	if strings.HasPrefix(extension, "#") {
		folderPath = path.Join(fileData.folderPath, fileData.fileName)
	} else if strings.Contains(extension, "#") {
		fileName := strings.Split(extension, "#")[0]
		folderPath = path.Join(fileData.folderPath, fileName)
	}
	folderPath, _ = url.PathUnescape(folderPath)
	return folderPath
}

func (fileData *FileData) File() *os.File {
	return fileData.file
}

func NewFileData(path string) (*FileData, error) {
	fileData := &FileData{}
	fileInfo, err := os.Stat(path)
	if err != nil {
		println("Failed to read file " + path + " " + err.Error())
		return nil, err
	}
	if fileInfo.IsDir() {
		fileData.folderPath = path
	} else {
		fileData.folderPath, _ = filepath.Split(path)
		fileData.fileName = filepath.Base(path)
		fileData.file, _ = os.Open(path)
		fileData.scanner = bufio.NewScanner(fileData.File())
		fileData.lineNumber = 1
	}
	return fileData, nil

}

func (fileData *FileData) Close() {
	defer fileData.file.Close()
}

func (fileData *FileData) ScanOneLine() (string, int) {
	if fileData.scanner.Scan() {
		lineText := fileData.scanner.Text()
		fileData.lineNumber++
		return lineText, fileData.lineNumber
	}
	return "", -1
}

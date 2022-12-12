package main

import (
	"bufio"
	"linkcheck/models"

	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var wg sync.WaitGroup

var result = models.Result{
	FilesLinksMap: map[string]*models.FileLink{},
}

func main() {
	readmeFiles := extractReadmeFiles()

	extractLinksFromReadmes(readmeFiles)
	wg.Wait()
	result.Print()
}

func extractLinksFromReadmes(files []string) {
	for _, filePath := range files {
		wg.Add(1)
		go handleFile(filePath)
	}
}

func handleFile(filePath string) {
	defer wg.Done()
	fileLinkData := models.FileLink{
		FilePath: filePath,
		Links:    []models.Link{},
	}
	result.AddNewFile(&fileLinkData)
	fileBytes, err := os.Open(filePath)
	defer fileBytes.Close()
	scanner := bufio.NewScanner(fileBytes)

	lineNumber := 1
	for scanner.Scan() {
		lineText := scanner.Text()
		findAndCheckLinksInLine(filePath, lineText, lineNumber)
		lineNumber++
	}
	if err != nil {
		println("Failed to read file " + filePath + " " + err.Error())
	}
}

func findAndCheckLinksInLine(filePath string, line string, lineNumber int) {
	var linksPaths []string
	linkRegex, _ := regexp.Compile("\\[.*\\]\\(.*\\)|https?:\\/\\/(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()@:%_\\+.~#?&//=]*)")
	linksPaths = append(linksPaths, linkRegex.FindAllString(line, -1)...)
	for index, linkPath := range linksPaths {
		if strings.Contains(linkPath, "](") {
			linkPath = strings.Split(linkPath, "](")[1]
		}
		linkPath = strings.Split(linkPath, ")")[0]
		linksPaths[index] = linkPath
	}
	for _, linkPath := range linksPaths {
		wg.Add(1)
		linkData := models.Link{
			LineNumber: lineNumber,
			Status:     0,
			Path:       linkPath,
		}
		go checkLink(filePath, linkData)
	}
}

func checkLink(filePath string, linkData models.Link) {
	defer wg.Done()
	switch {
	case strings.Contains(linkData.Path, "http"):
		linkData.LinkType = models.URL
		var err error
		resp, err := &http.Response{
			StatusCode: 200,
		}, nil //http.Get(linkData.path)
		if err != nil {
			println("Failed to get URL data with path " + linkData.Path + " and error " + err.Error())
			if strings.Contains(err.Error(), "timeout") {
				linkData.Status = 504
			}
		}
		if resp != nil {
			linkData.Status = resp.StatusCode
		}
	case strings.Contains(linkData.Path, "mailto:"):
		linkData.LinkType = models.Email
		mailRegex, _ := regexp.Compile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$")
		email := strings.Split(linkData.Path, ":")[0]
		if !mailRegex.MatchString(email) {
			linkData.Status = 400
			return
		}
		linkData.Status = 200
	default:
		linkData.LinkType = models.Folder
		folderPath, _ := filepath.Split(filePath)
		pathWithoutTitleLink := strings.Split(linkData.Path, "#")[0]
		folderPath = filepath.Join(folderPath, pathWithoutTitleLink)
		fileBytes, err := os.ReadFile(folderPath)
		if err != nil {
			linkData.Status = 400
			println("Failed to get link data with path " + linkData.Path + " and error " + err.Error())
			return
		}
		if strings.Contains(linkData.Path, "#") {
			fileData := string(fileBytes)
			if !fileContainsLink(linkData.Path, fileData) {
				linkData.Status = 400
				return
			}
		}
		linkData.Status = 200
	}
}

func fileContainsLink(titleLink string, fileText string) bool {
	titleLink = strings.Split(titleLink, "#")[1]
	title := strings.ReplaceAll(titleLink, "#", "")
	title = strings.ReplaceAll(title, "-", "( |-|)")
	readmeTitleRegex := "(?i)#( ?)" + title
	linkRegex, _ := regexp.Compile(readmeTitleRegex)
	if len(linkRegex.FindStringSubmatch(fileText)) > 0 {
		return true
	}
	return false
}

func extractReadmeFiles() []string {
	path := ""
	var readmeFiles []string

	if envPath := os.Getenv("PROJECT_PATH"); envPath != "" {
		path = envPath
	}
	err := filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
		if file.IsDir() && strings.Contains(file.Name(), "vendor") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
			path, _ = filepath.Abs(path)
			readmeFiles = append(readmeFiles, path)
		}
		return nil
	})

	if err != nil {
		println("Failed to get files with error " + err.Error())
	}
	return readmeFiles
}

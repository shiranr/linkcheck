package main

import (
	"linkcheck/models"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var lineHandler = models.GetInstance()

var result = models.Result{
	FilesLinksMap: map[string]*models.FileLink{},
}

// TODO add CMD.
// TODO make this a linter for megalinter.
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
		Links:    []*models.Link{},
	}
	result.AddNewFile(&fileLinkData)
	fileData, err := models.NewFileData(filePath)
	if err != nil {
		return
	}
	lineText, lineNumber := fileData.ScanOneLine()
	for lineNumber != -1 {
		linksPaths := lineHandler.FindAndCheckLinksInLine(lineText)
		for _, linkPath := range linksPaths {
			linkData := &models.Link{
				LineNumber: lineNumber,
				Status:     0,
				Path:       linkPath,
			}
			wg.Add(1)
			go checkLink(fileData, linkData)
		}
		lineText, lineNumber = fileData.ScanOneLine()
	}
}

func checkLink(fileData *models.FileData, linkData *models.Link) {
	defer wg.Done()
	switch {
	case strings.HasPrefix(linkData.Path, "http"):
		linkData.LinkType = models.URL
		var err error
		resp, err := httpRequest(linkData.Path)
		if err != nil {
			println("Failed to get URL data with path " + linkData.Path + " and error " + err.Error())
			if strings.Contains(err.Error(), "timeout") {
				linkData.Status = 504
			}
		}
		if resp != nil {
			linkData.Status = resp.StatusCode
		}
	case strings.HasPrefix(linkData.Path, "mailto:"):
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
		linkedFileEscapedFullPath := fileData.EscapedFullPath(linkData.Path)
		_, err := os.Stat(linkedFileEscapedFullPath)
		if err != nil {
			linkData.Status = 400
			println("Failed to get link data with path " + linkData.Path + " and error " + err.Error())
			return
		}
		if strings.Contains(linkData.Path, "#") {
			fileBytes, _ := os.ReadFile(linkedFileEscapedFullPath)
			fileData := string(fileBytes)
			if !fileContainsLink(linkData.Path, fileData) {
				linkData.Status = 400
				return
			}
		}
		linkData.Status = 200
	}
	result.Append(linkData, fileData.FullFilePath())
}

func httpRequest(link string) (*http.Response, error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("HEAD", link, nil)
	req.Header.Set("User-Agent", "Golang_Link_Check/1.0")
	resp, err := client.Do(req)
	for i := 0; i < 2 && ((resp == nil && err != nil) || (resp != nil && resp.StatusCode == 404 || resp.StatusCode == 403)); i++ {
		req, err = http.NewRequest("GET", link, nil)
		req.Header.Set("User-Agent", "Golang_Link_Check/1.0")
		resp, err = client.Do(req)
	}
	return resp, err
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
	path := "./"
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

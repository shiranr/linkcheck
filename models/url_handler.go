package models

import (
	"net/http"
	"strings"
	"time"
)

var (
	handler *urlHandler
)

type urlHandler struct {
	client *http.Client
}

func GetURLHandlerInstance() *urlHandler {
	if handler == nil {
		timeout := time.Duration(5 * time.Second)
		client := &http.Client{
			Timeout: timeout,
		}
		handler = &urlHandler{
			client: client,
		}
	}
	return handler
}

func (handler *urlHandler) Handle(linkData *Link) {
	linkData.LinkType = URL
	var err error
	resp, err := handler.httpRequest(linkData.Path)
	if err != nil {
		println("Failed to get URL data with path " + linkData.Path + " and error " + err.Error())
		if strings.Contains(err.Error(), "timeout") {
			linkData.Status = 504
		}
	}
	if resp != nil {
		defer resp.Body.Close()
		linkData.Status = resp.StatusCode
	}
}

func (handler *urlHandler) httpRequest(link string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", link, nil)
	req.Header.Set("User-Agent", "Golang_Link_Check/1.0")
	resp, err := handler.client.Do(req)
	for i := 0; i < 2 && ((resp == nil && err != nil) || (resp != nil && resp.StatusCode == 404 || resp.StatusCode == 403)); i++ {
		req, err = http.NewRequest("GET", link, nil)
		req.Header.Set("User-Agent", "Golang_Link_Check/1.0")
		resp, err = handler.client.Do(req)
	}
	return resp, err
}

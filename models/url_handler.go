package models

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
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
		// TODO make configurable
		timeout := time.Duration(30 * time.Second)
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.MaxIdleConns = 1000
		transport.IdleConnTimeout = timeout
		transport.DialContext = (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		}).DialContext
		transport.TLSHandshakeTimeout = timeout
		client := &http.Client{
			Timeout:   timeout,
			Transport: transport,
		}
		handler = &urlHandler{
			client: client,
		}
	}
	return handler
}

func (handler *urlHandler) Handle(linkData *Link) {
	linkData.LinkType = URL
	respStatus, err := handler.httpRequest(linkData.Path)
	if err != nil {
		log.WithFields(log.Fields{
			"link":  linkData.Path,
			"error": err,
		}).Error("Failed get URL data")
		if strings.Contains(err.Error(), "timeout") {
			linkData.Status = 504
		}
	}
	linkData.Status = respStatus
}

func (handler *urlHandler) httpRequest(link string) (int, error) {
	resp, err := handler.sendRequest("HEAD", link)
	if err != nil {
		return 0, err
	}
	for i := 0; i < 2 && ((resp == nil && err != nil) || (resp != nil && resp.StatusCode == 404 || resp.StatusCode == 403)); i++ {
		resp, err = handler.sendRequest("GET", link)
		if err != nil {
			return 0, err
		}
	}
	return resp.StatusCode, nil
}

func (handler *urlHandler) sendRequest(method string, link string) (*http.Response, error) {
	req, err := handler.createRequest(link, method)
	if err != nil {
		return nil, err
	}
	resp, err := handler.client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"link":  link,
			"error": err,
		}).Error("Failed to perform request")
		return nil, err
	}
	resp.Body.Close()
	return resp, nil
}

func (handler *urlHandler) createRequest(link string, method string) (*http.Request, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	req, err := http.NewRequest(method, link, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"link":   link,
			"method": method,
			"error":  err,
		}).Error("Failed to create timeout for request")
		return nil, err
	}
	req.WithContext(ctx)
	req.Close = true
	req.Header.Set("User-Agent", "Golang_Link_Check/1.0")
	return req, err
}

package models

import (
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	handler *urlHandler
)

type urlHandler struct {
}

// GetURLHandlerInstance - get instance of URL handler, handle links with http in them (singleton)
func GetURLHandlerInstance() LinkHandlerInterface {
	if handler == nil {
		handler = &urlHandler{}
	}
	return handler
}

// Handle - using scrap lib, check the link status
func (handler *urlHandler) Handle(linkPath string) int {
	respStatus, err := handler.scrap(linkPath)
	for i := 0; i < 3 && (err != nil || !handler.respStatusOK(respStatus)); i++ {
		log.WithFields(log.Fields{
			"link":       linkPath,
			"error":      err,
			"respStatus": respStatus,
		}).Error("Failed get URL data")
		time.Sleep(500 * time.Millisecond)
		errLower := strings.ToLower(err.Error())
		respStatus, err = handler.scrap(linkPath)
		if err == nil && handler.respStatusOK(respStatus) {
			return respStatus
		}
		errLower = strings.ToLower(err.Error())
		if strings.Contains(errLower, "not found") {
			return 404
		}
		if strings.Contains(errLower, "forbidden") {
			return 403
		}
	}
	return respStatus
}

func (handler *urlHandler) respStatusOK(restStatus int) bool {
	return restStatus >= 200 && restStatus < 300 || restStatus >= 400 && restStatus < 500
}

func (handler *urlHandler) scrap(linkPath string) (int, error) {
	c := colly.NewCollector()
	respStatus := 0
	c.OnResponse(func(resp *colly.Response) {
		respStatus = resp.StatusCode
	})
	return respStatus, c.Visit(linkPath)
}

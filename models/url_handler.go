package models

import (
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"strings"
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
	respStatus, err := handler.scrap(linkPath, true)
	for i := 0; i < 2 && err != nil; i++ {
		errLower := strings.ToLower(err.Error())
		respStatus, err = handler.scrap(linkPath, false)
		if err == nil {
			return respStatus
		}
		errLower = strings.ToLower(err.Error())
		if strings.Contains(errLower, "not found") {
			return 404
		}
		if strings.Contains(errLower, "forbidden") {
			return 403
		}
		if strings.Contains(errLower, "timeout") {
			return 504
		}
		log.WithFields(log.Fields{
			"link":  linkPath,
			"error": err,
		}).Error("Failed get URL data")
	}
	return respStatus
}

func (handler *urlHandler) respStatusOK(restStatus int) bool {
	return restStatus >= 200 && restStatus < 300 || restStatus >= 400 && restStatus < 500
}

func (handler *urlHandler) scrap(linkPath string, checkHead bool) (int, error) {
	c := colly.NewCollector()
	c.CheckHead = checkHead
	respStatus := 0
	c.OnResponse(func(resp *colly.Response) {
		respStatus = resp.StatusCode
	})
	return respStatus, c.Visit(linkPath)
}

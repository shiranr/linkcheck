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
	collector *colly.Collector
}

// GetURLHandlerInstance - get instance of URL handler, handle links with http in them (singleton)
func GetURLHandlerInstance() LinkHandlerInterface {
	if handler == nil {
		handler = &urlHandler{
			colly.NewCollector(),
		}
	}
	return handler
}

// Handle - using scrap lib, check the link status
func (handler *urlHandler) Handle(linkPath string) int {
	respStatus, err := handler.scrap(linkPath)
	for i := 0; i < 2 && err != nil; i++ {
		errLower := strings.ToLower(err.Error())
		if strings.Contains(errLower, "eof") || strings.Contains(errLower, "timeout") {
			respStatus, err = handler.scrap(linkPath)
			if err == nil {
				return respStatus
			}
			errLower = strings.ToLower(err.Error())
		}
		if strings.Contains(errLower, "not found") {
			return 404
		}
		if strings.Contains(errLower, "forbidden") {
			return 403
		}
		log.WithFields(log.Fields{
			"link":  linkPath,
			"error": err,
		}).Error("Failed get URL data")
	}
	return respStatus
}

func (handler *urlHandler) scrap(linkPath string) (int, error) {
	respStatus := 0
	handler.collector.OnResponse(func(resp *colly.Response) {
		respStatus = resp.StatusCode
	})
	return respStatus, handler.collector.Visit(linkPath)
}

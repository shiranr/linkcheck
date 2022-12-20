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

func GetURLHandlerInstance() *urlHandler {
	if handler == nil {
		handler = &urlHandler{}
	}
	return handler
}

func (handler *urlHandler) Handle(linkPath string) int {
	respStatus, err := handler.scrap(linkPath)
	if err != nil {
		log.WithFields(log.Fields{
			"link":  linkPath,
			"error": err,
		}).Error("Failed get URL data")
		if strings.Contains(err.Error(), "timeout") {
			return 504
		}
	}
	return respStatus
}

func (handler *urlHandler) scrap(linkPath string) (int, error) {
	c := colly.NewCollector()
	respStatus := 0
	c.OnResponse(func(resp *colly.Response) {
		respStatus = resp.StatusCode
	})
	return respStatus, c.Visit(linkPath)
}

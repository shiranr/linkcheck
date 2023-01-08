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

func GetURLHandlerInstance() LinkHandlerInterface {
	if handler == nil {
		handler = &urlHandler{}
	}
	return handler
}

func (handler *urlHandler) Handle(linkPath string) int {
	respStatus, err := handler.scrap(linkPath)
	for i := 0; i < 2 && err != nil; i++ {
		errLower := strings.ToLower(err.Error())
		if strings.Contains(errLower, "eof") {
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

func (handler *urlHandler) scrap(linkPath string) (int, error) {
	c := colly.NewCollector()
	respStatus := 0
	c.OnResponse(func(resp *colly.Response) {
		respStatus = resp.StatusCode
	})
	return respStatus, c.Visit(linkPath)
}

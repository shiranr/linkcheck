package models

import (
	"github.com/gocolly/colly"
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
	if err != nil {
		errLower := strings.ToLower(err.Error())
		if strings.Contains(errLower, "not found") {
			return 404
		}
		if strings.Contains(errLower, "forbidden") {
			return 403
		}
		if strings.Contains(errLower, "timeout") {
			return 504
		}
	}
	return respStatus
}

func (handler *urlHandler) scrap(linkPath string, headFirst bool) (int, error) {
	var err error
	c := colly.NewCollector()
	respStatus := 0
	c.OnResponse(func(resp *colly.Response) {
		respStatus = resp.StatusCode
	})
	if headFirst {
		c.CheckHead = true
		err = c.Visit(linkPath)
	}
	if err != nil || respStatus == 0 {
		c.CheckHead = false
		err = c.Visit(linkPath)
	}
	return respStatus, err
}

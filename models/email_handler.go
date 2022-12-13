package models

import (
	"regexp"
	"strings"
)

var (
	emailRegex, _ = regexp.Compile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$")
	mailHandler   *emailHandler
)

type emailHandler struct {
}

func GetEmailHandlerInstance() *emailHandler {
	if mailHandler == nil {
		mailHandler = &emailHandler{}
	}
	return mailHandler
}

func (handler *emailHandler) Handle(linkData *Link) {
	linkData.LinkType = Email
	email := strings.Split(linkData.Path, ":")[0]
	if !emailRegex.MatchString(email) {
		linkData.Status = 400
		return
	}
	linkData.Status = 200
}

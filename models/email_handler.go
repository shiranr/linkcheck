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

func (handler *emailHandler) Handle(linkPath string) int {
	email := strings.Split(linkPath, ":")[0]
	if !emailRegex.MatchString(email) {
		return 400
	}
	return 200
}

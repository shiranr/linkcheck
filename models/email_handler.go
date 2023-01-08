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

// GetEmailHandlerInstance Get email handler instance
func GetEmailHandlerInstance() LinkHandlerInterface {
	if mailHandler == nil {
		mailHandler = &emailHandler{}
	}
	return mailHandler
}

func (handler *emailHandler) Handle(linkPath string) int {
	email := strings.Split(linkPath, ":")[1]
	if !emailRegex.MatchString(email) {
		return 400
	}
	return 200
}

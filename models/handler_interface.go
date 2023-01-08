package models

type LinkHandlerInterface interface {
	Handle(linkPath string) int
}

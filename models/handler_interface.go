package models

// LinkHandlerInterface - interface to define different link handlers
type LinkHandlerInterface interface {
	Handle(linkPath string) int
}

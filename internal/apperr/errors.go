package apperr

import "errors"

var (
	ErrorLinkNotFound       = errors.New("link not found")
	ErrorEmptyCreateRequest = errors.New("empty create request")
)

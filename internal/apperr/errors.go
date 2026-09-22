package apperr

import "errors"

var (
	ErrorLinkNotFound         = errors.New("link not found")
	ErrorInvalidCreateRequest = errors.New("invalid create request")
)

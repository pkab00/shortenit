package apperr

import "errors"

var (
	ErrorLinkNotFound         = errors.New("link not found")
	ErrorInvalidCreateRequest = errors.New("invalid create request")
	ErrorInvalidCustomCode    = errors.New("invalid custom code")
	ErrorCustomCodeInUse      = errors.New("custom code is already in use")
)

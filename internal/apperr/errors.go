package apperr

import "errors"

var (
	ErrorLinkNotFound         = errors.New("link not found")
	ErrorInvalidCreateRequest = errors.New("invalid create request")
	ErrorCustomCodeTooLong    = errors.New("custom code is too long")
	ErrorCustomCodeInUse      = errors.New("custom code is already in use")
)

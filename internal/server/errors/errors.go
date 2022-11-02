package errors

import "errors"

var (
	ErrInvalidBody      = errors.New("invalid request body")
	ErrCreatingInstance = errors.New("error creating instance")
)

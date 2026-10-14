package model

import "errors"

var (
	ErrNotFound = errors.New("not found error")
	ErrConflict = errors.New("conflict error")
	ErrValidation = errors.New("validation error")
	ErrMissingConfig = errors.New("missing environment variable")
)
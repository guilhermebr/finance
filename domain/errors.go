package domain

import "errors"

var (
	ErrConflict            = errors.New("data conflict")
	ErrNotFound            = errors.New("not found")
	ErrMalformedParameters = errors.New("malformed parameters")
	ErrForbidden           = errors.New("forbidden")
	ErrDuplicateKey        = errors.New("duplicate key")
)

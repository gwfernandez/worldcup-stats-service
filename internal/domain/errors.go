package domain

import "errors"

var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")
	// ErrDuplicateKey is returned when a unique constraint is violated
	ErrDuplicateKey = errors.New("resource already exists")
	// ErrInvalidInput is returned when the input data is invalid
	ErrInvalidInput = errors.New("invalid input data")
)

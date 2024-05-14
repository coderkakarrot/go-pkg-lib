package rest

import (
	"errors"
)

var (
	// ErrInternalServer is the error returned when an internal server error
	// occurs.
	ErrInternalServer = errors.New("internal server error")

	// ErrMissingContext is the error returned when the api value is missing
	// from the context.
	ErrMissingContext = errors.New("api value missing from context")
)

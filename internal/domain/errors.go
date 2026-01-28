package domain

import "errors"

var (
	// Standard Errors
	ErrNotFound           = errors.New("resource not found")
	ErrConflict           = errors.New("resource already exists") // duplicate key
	ErrInternal           = errors.New("internal system error")
	ErrInvalidInput       = errors.New("invalid input parameter")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("access forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

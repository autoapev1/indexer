package auth

import (
	"errors"
	"net/http"
)

type Provider interface {
	// Authenticate authenticates a request.
	Authenticate(r *http.Request) error
	Register() (key string, err error)
	UpdateUsage(key string, usageDelta KeyUsage) error
}

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidKey   = errors.New("invalid key")
)

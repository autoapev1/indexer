package auth

import (
	"errors"
	"net/http"
)

type AuthProvider string

const (
	AuthProviderNoAuth AuthProvider = "noauth"
	AuthProviderMemory AuthProvider = "memory"
	AuthProviderSql    AuthProvider = "sql"
)

func ToProvider(s string) AuthProvider {
	switch s {
	case "noauth":
		return AuthProviderNoAuth
	case "memory":
		return AuthProviderMemory
	case "sql":
		return AuthProviderSql
	default:
		return AuthProviderNoAuth
	}
}

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

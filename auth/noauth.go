package auth

import "net/http"

type NoAuthProvider struct{}

func NewNoAuthProvider() *NoAuthProvider {
	return &NoAuthProvider{}
}

func (a *NoAuthProvider) Authenticate(r *http.Request) error {
	return nil
}

func (a *NoAuthProvider) Register() (string, error) {
	return "", nil
}

func (a *NoAuthProvider) UpdateUsage(key string, usageDelta KeyUsage) error {
	return nil
}

// ensure NoAuth implements Provider
var _ Provider = (*NoAuthProvider)(nil)

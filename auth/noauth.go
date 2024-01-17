package auth

import "net/http"

type NoAuth struct{}

func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

func (a *NoAuth) Authenticate(r *http.Request) error {
	return nil
}

func (a *NoAuth) Register() (string, error) {
	return "", nil
}

func (a *NoAuth) UpdateUsage(key string, usageDelta *KeyUsage) error {
	return nil
}

// ensure NoAuth implements Provider
var _ Provider = (*NoAuth)(nil)

package auth

import "net/http"

type NoAuthProvider struct{}

func NewNoAuthProvider() *NoAuthProvider {
	return &NoAuthProvider{}
}

// with no auth, highest auth level is defualt
func (a *NoAuthProvider) Authenticate(r *http.Request) (AuthLevel, error) {
	return AuthLevelMaster, nil
}

func (a *NoAuthProvider) Register() (string, error) {
	return "", nil
}

func (a *NoAuthProvider) UpdateUsage(key string, usageDelta KeyUsage) error {
	return nil
}

// ensure NoAuth implements Provider
var _ Provider = (*NoAuthProvider)(nil)

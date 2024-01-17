package auth

import (
	"sync"

	"github.com/savsgio/gotils/nocopy"
)

type MemoryProvider struct {
	nocopy.NoCopy
	// lock protects the keys map.
	lock sync.RWMutex
	// keys is a map of API keys to their corresponding secret.
	Keys map[string]*KeyUsage

	KeyType KeyType
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{
		Keys: make(map[string]*KeyUsage),
	}
}

func (a *MemoryProvider) Authenticate(key string) error {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if _, ok := a.Keys[key]; !ok {
		return ErrInvalidKey
	}

	return nil
}

func (a *MemoryProvider) Register() (string, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	key, err := GenerateKey(a.KeyType)
	if err != nil {
		return "", err
	}

	a.Keys[key] = &KeyUsage{}

	return key, nil
}

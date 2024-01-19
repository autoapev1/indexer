package auth

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/autoapev1/indexer/config"
	"github.com/savsgio/gotils/nocopy"
)

type memoryKey struct {
	Iat         int64            `json:"iat"`
	Exp         int64            `json:"exp"`
	LastIP      string           `json:"last_ip"`
	LastAccess  int64            `json:"last_access"`
	CallCount   int64            `json:"call_count"`
	MethodUsage map[string]int64 `json:"method_usage"`
}

type MemoryProvider struct {
	nocopy.NoCopy
	// lock protects the keys map.
	lock sync.RWMutex
	// keys is a map of API keys to their corresponding secret.
	Keys map[string]*memoryKey

	keyType KeyType

	defaultExpiry time.Duration
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{
		Keys:          make(map[string]*memoryKey),
		keyType:       KeyTypeHex64,
		defaultExpiry: time.Hour * 24 * 30 * 3, // 3 months
	}
}

func (a *MemoryProvider) WithKeyType(keyType KeyType) *MemoryProvider {
	a.keyType = keyType
	return a
}

func (a *MemoryProvider) WithDefaultExpiry(defaultExpiry time.Duration) *MemoryProvider {
	if defaultExpiry < 0 {
		slog.Warn("invalid default expiry, using default value", "default_expiry", defaultExpiry.String())
		return a
	}

	a.defaultExpiry = defaultExpiry
	return a
}

func (a *MemoryProvider) Authenticate(r *http.Request) (AuthLevel, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	key := r.Header.Get("Authentication")
	key = strings.TrimPrefix(key, "Bearer ")

	// check master
	master := config.Get().API.AuthMasterKey
	if master != "" && key == master {
		return AuthLevelMaster, nil
	}

	if _, ok := a.Keys[key]; !ok {
		return AuthLevelUnauthorized, ErrUnauthorized
	}

	return AuthLevelBasic, nil
}

func (a *MemoryProvider) Register() (string, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	key, err := GenerateKey(a.keyType)
	if err != nil {
		return "", err
	}

	a.Keys[key] = &memoryKey{}

	return key, nil
}

func (a *MemoryProvider) UpdateUsage(key string, usageDelta KeyUsage) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, ok := a.Keys[key]; !ok {
		return ErrInvalidKey
	}

	a.Keys[key].CallCount += usageDelta.Requests
	a.Keys[key].LastIP = usageDelta.IP
	a.Keys[key].LastAccess = usageDelta.AccessedAt
	for k, v := range usageDelta.MethodUsage {
		a.Keys[key].MethodUsage[k] += v
	}

	return nil
}

// ensure MemoryProvider implements Provider
var _ Provider = (*MemoryProvider)(nil)

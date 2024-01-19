package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

// context key for middleware
type CtxAuthKey int

const AuthKey CtxAuthKey = 0

// auth levels
type AuthLevel int

const (
	AuthLevelUnauthorized AuthLevel = iota
	AuthLevelBasic
	AuthLevelMaster
)

func IsValidAuthLevel(lvl AuthLevel) bool {
	switch lvl {
	case AuthLevelUnauthorized, AuthLevelBasic, AuthLevelMaster:
		return true
	default:
		return false
	}
}

// key types
type KeyType int

const (
	KeyTypeUUID KeyType = iota
	KeyTypeHex16
	KeyTypeHex32
	KeyTypeHex64
	KeyTypeHex128
	KeyTypeHex256
)

// string key types
const (
	KeyTypeUUIDString   = "uuid"
	KeyTypeHex16String  = "hex16"
	KeyTypeHex32String  = "hex32"
	KeyTypeHex64String  = "hex64"
	KeyTypeHex128String = "hex128"
	KeyTypeHex256String = "hex256"
)

// take string of key type to int KeyType
func ToKeyType(s string) KeyType {
	switch s {
	case KeyTypeUUIDString:
		return KeyTypeUUID
	case KeyTypeHex16String:
		return KeyTypeHex16
	case KeyTypeHex32String:
		return KeyTypeHex32
	case KeyTypeHex64String:
		return KeyTypeHex64
	case KeyTypeHex128String:
		return KeyTypeHex128
	case KeyTypeHex256String:
		return KeyTypeHex256
	default:
		slog.Warn("invalid key type", "KeyType", s)
		return KeyTypeHex64
	}
}

func GenerateKey(keyType KeyType) (string, error) {
	switch keyType {
	case KeyTypeUUID:
		return uuid.New().String(), nil
	case KeyTypeHex16:
		return generateRandomHex(8) // 16 hex characters
	case KeyTypeHex32:
		return generateRandomHex(16) // 32 hex characters
	case KeyTypeHex64:
		return generateRandomHex(32) // 64 hex characters
	case KeyTypeHex128:
		return generateRandomHex(64) // 128 hex characters
	case KeyTypeHex256:
		return generateRandomHex(128) // 256 hex characters
	default:
		return "", errors.New("invalid key type")
	}
}

// generateRandomHex generates a random hexadecimal string of the specified byte length.
func generateRandomHex(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type KeyUsage struct {
	IP          string
	AccessedAt  int64
	Requests    int64
	MethodUsage map[string]int64
}

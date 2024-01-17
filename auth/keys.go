package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
)

type KeyType int

const (
	KeyTypeUUID KeyType = iota
	KeyTypeHex16
	KeyTypeHex32
	KeyTypeHex64
	KeyTypeHex128
	KeyTypeHex256
)

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

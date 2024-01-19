package api

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/autoapev1/indexer/auth"
)

type JRPCRequest struct {
	ID      string          `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type JRPCResponse struct {
	ID      string          `json:"id,omitempty"`
	JSONRPC string          `json:"jsonrpc,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JRPCError      `json:"error,omitempty"`
}

type Response interface{}

type JRPCError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type MethodPrefix string

const (
	MethodInvalid MethodPrefix = ""
	MethodIdx     MethodPrefix = "idx_"
	MethodAuth    MethodPrefix = "auth_"
)

func getMethodPrefix(s string) MethodPrefix {
	if !strings.Contains(s, "_") {
		return MethodInvalid
	}
	m := strings.Split(s, "_")[0]

	switch m {
	case "idx":
		return MethodIdx
	case "auth":
		return MethodAuth
	default:
		return MethodInvalid
	}
}

func hasAccess(methodPrefix MethodPrefix, authlvl auth.AuthLevel) bool {
	switch methodPrefix {
	case MethodIdx:
		return authlvl >= auth.AuthLevelBasic
	case MethodAuth:
		return authlvl >= auth.AuthLevelMaster
	default:
		slog.Warn("invalid method prefix", "method_prefix", methodPrefix, "auth_level", authlvl)
		return false
	}
}

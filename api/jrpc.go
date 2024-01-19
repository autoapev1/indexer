package api

import "encoding/json"

type JRPCRequest struct {
	ID      int64           `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type JRPCResponse struct {
	ID      int64           `json:"id,omitempty"`
	JSONRPC string          `json:"jsonrpc,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JRPCError      `json:"error,omitempty"`
}

type JRPCError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

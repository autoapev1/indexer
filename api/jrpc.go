package api

import "encoding/json"

type JRPCRequest struct {
	ID      int64           `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type JRPCResponse struct {
	ID      int64           `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *JRPCError      `json:"error"`
}

type JRPCError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

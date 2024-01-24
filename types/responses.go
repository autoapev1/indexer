package types

type JRPCError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type GetBlockNumberResponse struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Result map[int64]int64 `json:"result,omitempty"`
	Error  *JRPCError      `json:"error,omitempty"`
}

type GetHeightsResponse struct {
	ID     string     `json:"id"`
	Method string     `json:"method"`
	Result *Heights   `json:"result,omitempty"`
	Error  *JRPCError `json:"error,omitempty"`
}

type GetChainsResponse struct {
	ID     string     `json:"id"`
	Method string     `json:"method"`
	Result []Chain    `json:"result,omitempty"`
	Error  *JRPCError `json:"error,omitempty"`
}

type GetBlockTimestampsResponse struct {
	ID     string            `json:"id"`
	Method string            `json:"method"`
	Result []*BlockTimestamp `json:"result,omitempty"`
	Error  *JRPCError        `json:"error,omitempty"`
}

type GetBlockAtTimestampResponse struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Result *BlockTimestamp `json:"result,omitempty"`
	Error  *JRPCError      `json:"error,omitempty"`
}

type FindTokensResponse struct {
	ID     string     `json:"id"`
	Method string     `json:"method"`
	Result []*Token   `json:"result,omitempty"`
	Error  *JRPCError `json:"error,omitempty"`
}

type FindPairsResponse struct {
	ID     string     `json:"id"`
	Method string     `json:"method"`
	Result []*Pair    `json:"result,omitempty"`
	Error  *JRPCError `json:"error,omitempty"`
}

type GetTokenCountResponse struct {
	ID     string     `json:"id"`
	Method string     `json:"method"`
	Result int64      `json:"result,omitempty"`
	Error  *JRPCError `json:"error,omitempty"`
}

type GetPairCountResponse struct {
	ID     string     `json:"id"`
	Method string     `json:"method"`
	Result int64      `json:"result,omitempty"`
	Error  *JRPCError `json:"error,omitempty"`
}

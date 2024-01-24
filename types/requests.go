package types

type GetBlockNumberRequest struct{}

type GetChainsRequest struct{}

type GetHeightsRequest struct {
	ChainID *int64 `json:"chain_id"`
}

type GetBlockTimestampsRequest struct {
	ChainID   *int64 `json:"chain_id"`
	FromBlock *int64 `json:"from_block"`
	ToBlock   *int64 `json:"to_block"`
}

type GetBlockAtTimestampRequest struct {
	ChainID   *int64 `json:"chain_id"`
	Timestamp *int64 `json:"timestamp"`
}

type FindTokensRequest struct {
	ChainID *int64        `json:"chain_id"`
	Filter  *TokenFilter  `json:"filter,omitempty"`
	Options *TokenOptions `json:"options,omitempty"`
}

type GetTokenCountRequest struct {
	ChainID *int64 `json:"chain_id"`
}

type FindPairsRequest struct {
	ChainID *int64       `json:"chain_id"`
	Filter  *PairFilter  `json:"filter,omitempty"`
	Options *PairOptions `json:"options,omitempty"`
}

type GetPairCountRequest struct {
	ChainID *int64 `json:"chain_id"`
}

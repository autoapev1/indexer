package types

type GetBlockNumberRequest struct{}

type GetChainsRequest struct{}

type GetBlockTimestampsRequest struct {
	ChainID   int64 `json:"chain_id"`
	FromBlock int64 `json:"from_block"`
	ToBlock   int64 `json:"to_block"`
}

type GetBlockAtTimestampRequest struct {
	ChainID   int64 `json:"chain_id"`
	Timestamp int64 `json:"timestamp"`
}

type FindTokensRequest struct {
	ChainID int64        `json:"chain_id"`
	Filter  TokenFilter  `json:"filter"`
	Options TokenOptions `json:"options"`
}

type GetTokenCountRequest struct {
	ChainID int64 `json:"chain_id"`
}

type FindPairsRequest struct {
	ChainID int64       `json:"chain_id"`
	Filter  PairFilter  `json:"filter"`
	Options PairOptions `json:"options"`
}

type GetPairCountRequest struct {
	ChainID int64 `json:"chain_id"`
}

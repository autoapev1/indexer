package types

type TokenFilter struct {
	Address   *string `json:"address,omitempty"`
	Creator   *string `json:"creator,omitempty"`
	FromBlock *int64  `json:"from_block,omitempty"`
	ToBlock   *int64  `json:"to_block,omitempty"`
	Name      *string `json:"name,omitempty"`
	Symbol    *string `json:"symbol,omitempty"`
	Decimals  *uint8  `json:"decimals,omitempty"`
	Fuzzy     bool    `json:"fuzzy"`
}

type PairFilter struct {
	Token0Address *string `json:"token0_address,omitempty"`
	Token1Address *string `json:"token1_address,omitempty"`
	PoolAddress   *string `json:"pool_address,omitempty"`
	FromBlock     *int64  `json:"from_block,omitempty"`
	ToBlock       *int64  `json:"to_block,omitempty"`
	Fee           *int64  `json:"fee,omitempty"`
	TickSpacing   *int64  `json:"tick_spacing,omitempty"`
	Hash          *string `json:"hash,omitempty"`
	PoolType      *uint8  `json:"pool_type,omitempty"`
	Fuzzy         bool    `json:"fuzzy"`
}

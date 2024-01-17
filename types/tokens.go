package types

type TokenInfo struct {
	Address        string `json:"address" bun:",pk"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	Decimals       int    `json:"decimals"`
	Creator        string `json:"creator"`
	CreatedAtBlock int    `json:"created_at"`
	ChainID        int    `json:"chain_id"`
}

type BlockTimestamp struct {
	Block     int `json:"b" bun:",pk"`
	Timestamp int `json:"t"`
}

type PairInfo struct {
	Token0      string `json:"token0" bun:",pk"`
	Token1      string `json:"token1" bun:",pk"`
	Fee         int    `json:"fee"`
	TickSpacing int    `json:"tick_spacing"`
	Pool        string `json:"pool"`
	PoolType    int    `json:"pool_type"`
	CreatedAt   int    `json:"created_at"`
	Hash        string `json:"hash"`
	ChainID     int    `json:"chain_id"`
}

type OHLC struct {
	TS uint32  `json:"ts"`  // timestamp
	US float32 `json:"usd"` // usd price
	O  float32 `json:"o"`   // open price
	H  float32 `json:"h"`   // high price
	L  float32 `json:"l"`   // low price
	C  float32 `json:"c"`   // close price
	BV float32 `json:"bv"`  // buy volume usd
	SV float32 `json:"sv"`  // sell volume usd
	NB uint32  `json:"nb"`  // number of buy
	NS uint32  `json:"ns"`  // number of sells
}

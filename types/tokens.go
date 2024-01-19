package types

type Token struct {
	Address        string `json:"address" bun:",pk"`
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	Decimals       uint8  `json:"decimals"`
	Creator        string `json:"creator"`
	CreatedAtBlock int64  `json:"created_at"`
	ChainID        int16  `json:"chain_id"`
}

type BlockTimestamp struct {
	Block     int64 `json:"block" bun:",pk"`
	Timestamp int64 `json:"timestamp"`
}

type Pair struct {
	Token0Address string `json:"token0_address"`
	Token1Address string `json:"token1_address"`
	Fee           int64  `json:"fee"`
	TickSpacing   int64  `json:"tick_spacing"`
	PoolAddress   string `json:"pool_address"`
	PoolType      uint8  `json:"pool_type"`
	CreatedAt     int64  `json:"created_at"`
	Hash          string `json:"hash" bun:",pk"`
	ChainID       int16  `json:"chain_id"`
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

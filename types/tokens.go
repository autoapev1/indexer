package types

import (
	"strings"

	"github.com/uptrace/bun"
)

type Token struct {
	bun.BaseModel `bun:"table:tokens,alias:tokens" json:"-"`
	Address       string `json:"address" bun:",pk,type:varchar(42),unique"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	Decimals      uint8  `json:"decimals"`
	Creator       string `json:"creator" bun:",type:varchar(42),default:'0x0000000000000000000000000000000000000000'"`
	CreatedAt     int64  `json:"created_at"`
	CreationHash  string `json:"creation_hash" bun:",type:varchar(66),default:'0x0000000000000000000000000000000000000000000000000000000000000000'"`
	ChainID       int16  `json:"chain_id"`
}

func (p *Token) Lower() {
	if p.Name == "_Unknown" {
		p.Name = "unknown"
	}
	if p.Symbol == "_Unknown" {
		p.Symbol = "unknown"
	}

	p.Address = strings.ToLower(p.Address)
	p.Creator = strings.ToLower(p.Creator)
	p.CreationHash = strings.ToLower(p.CreationHash)
}

type BlockTimestamp struct {
	bun.BaseModel `bun:"table:block_timestamps,alias:block_timestamps" json:"-"`
	Block         int64 `json:"block" bun:",pk,notnull,unique"`
	Timestamp     int64 `json:"timestamp" bun:",notnull,default:0"`
}

type Creator struct {
	Hash    string `json:"hash"`
	Creator string `json:"creator"`
}

type BlockNumber struct {
	Number string `json:"blockNumber"`
}

type Pair struct {
	bun.BaseModel `bun:"table:pairs,alias:pairs" json:"-"`
	Token0Address string `json:"token0_address" bun:",type:varchar(42),notnull"`
	Token1Address string `json:"token1_address" bun:",type:varchar(42),notnull"`
	Fee           int64  `json:"fee" bun:",notnull,default:0"`
	TickSpacing   int64  `json:"tick_spacing" bun:",notnull,default:0"`
	PoolAddress   string `json:"pool_address" bun:",notnull,type:varchar(42),unique"`
	PoolType      uint8  `json:"pool_type" bun:",notnull,default:0"`
	CreatedAt     int64  `json:"created_at"`
	Hash          string `json:"hash" bun:",pk,type:varchar(66)"`
	ChainID       int16  `json:"chain_id"`
}

func (p *Pair) Lower() {
	p.Token0Address = strings.ToLower(p.Token0Address)
	p.Token1Address = strings.ToLower(p.Token1Address)
	p.PoolAddress = strings.ToLower(p.PoolAddress)
	p.Hash = strings.ToLower(p.Hash)
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

type Heights struct {
	Blocks int64
	Tokens int64
	Pairs  int64
}

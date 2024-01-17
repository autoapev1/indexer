# Evm Chain Indexer

## Overview

Indexer for EVM-based chains.
Currently supports Ethereum and Binance Smart Chain, using Uniswap and PancakeSwap as the DEXes.

### What does it Index?

- BlockTimestamps
- Token Info
- Pair Info
- Wallet Balances TODO
- Token Holders TODO
- Chart Data TODO

### How does it get data?

There are two ways to seed the database:

1. From scratch, using a private archive node. It will take a few days to build the indexes but is the cheapest option.
2. From a database dump. You can download a database dump from [TODO] and import it into your database. The database dump is updated weekly, so there will be a delay of up to 7 days which will need to be synced on the first run. This is the fastest option, but access to the database dump is not free.

Once the database is seeded, you will need access to any regular node to get each new block as it is mined. The database is updated in real-time, so as long as you keep the indexer running, it will stay up to date.

If the indexer is stopped for a while, it will need to catch up on the missed blocks. This can take a while, depending on how long it was stopped. Depending on how long it was stopped, you may need an archive node to access the data needed to resync, if this is the case you should consider using the database dump instead.

## Usage

### Requirements

- You will require a private (erigon) archive node to create the database from scratch, alternatively, you can download a database dump from [TODO]
- You will require a database to store the data. Currently, only Postgres is supported, but `storage.Store` is an interface so it's easy to add more.
  It is recommended to use NVMe SSD for storage, You will need around:
  - 7GB of disk space for BlockTimestamps (ETH + BSC)
  - 1GB of disk space for Tokens and Pairs (ETH + BSC)
  - 400GB of disk space for Wallet Balances (ETH + BSC) (estimate)
  - 400GB of disk space for Token Holders (ETH + BSC) (estimate)
  - 400GB of disk space for Chart Data (ETH + BSC) (estimate)

### Configuration

You will need to create a `config.toml` file in the root directory of the project. You can use [config.example.toml](config.example.toml) as a template.
Here is an example config for a local setup:

```toml
ethNodeAddr = "http://127.0.0.1:8500"
bscNodeAddr = "http://127.0.0.1:8400"

[tokens]
batchSize        = 10
batchConcurrency = 2

[pairs]
batchSize        = 10
batchConcurrency = 2
blockRange       = 200

[storage]
driver = "postgres"

[postgres]
user     = "postgres"
password = "postgres"
host     = "localhost"
port     = "5432"
sslmode  = "disable"

[api]
host    = "localhost"
port    = 8080
useAuth = true
apiKey  = "my-api-key"

```

### Running

```bash
go run ./cmd/api/main.go --config config.toml
```

### API

The API is JSON-RPC 2.0 compliant and is served on port 8080 by default.
The available methods are:

- `idx_getBlockTimestamps` - Get the timestamp for a range of block numbers
- `idx_getBlockAtTimestamp` - Get the block number at a timestamp

- `idx_getTokenByAddress` - Get info about a token by address
- `idx_getTokensByCreator` - Get info about tokens created by an address
- `idx_getTokensInBlock` - Get info about tokens created in a block range
- `idx_findTokens` - Find tokens by using find params

- `idx_getPairByAddress` - Get info about a pair by address
- `idx_getPairsByToken` - Get info about pairs containing a token
- `idx_getPairsInBlock` - Get info about pairs created in a block range
- `idx_findPairs` - Find pairs by using find params

### Indexed Types

#### BlockTimestamps

```go
type BlockTimestamp struct {
	BlockNumber uint64
	Timestamp   uint64
}
```

#### Token Info

```go
type Token struct {
	Address        string
	Name           string
	Symbol         string
	Decimals       uint8
	Creator        string
	CreatedAtBlock int64
	ChainID        int16
}
```

#### Pair Info

```go
type Pair struct {
	Token0Address string
	Token1Address string
	Fee           int64
	TickSpacing   int64
	PoolAddress   string
	PoolType      uint8
	CreatedAt     int64
	Hash          string
	ChainID       int16
}
```

#### OHLCVT Chart Data

```go
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
```

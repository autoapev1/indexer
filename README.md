# Evm Chain Indexer

## Overview

Indexer for EVM-based chains.
Currently supports Ethereum and Binance Smart Chain, using Uniswap and PancakeSwap as the DEXes.

### What does it Index?

- [x] BlockTimestamps
- [x] Token Info
- [x] Pair Info
- [ ] Wallet Balances
- [ ] Token Holders
- [ ] Liquidity Token Holders
- [ ] Chart Data

### How does it get data?

There are two ways to seed the database:

1. From scratch, using a private archive node. It will take a few days to build the indexes but is the cheapest option.
2. From a database dump. You can download a database dump from [TODO] and import it into your database. The database dump is updated weekly, so there will be a delay of up to 7 days which will need to be synced on the first run. This is the fastest option, but access to the database dump is not free.

Once the database is seeded, you will need access to any regular node to get each new block as it is mined. The database is updated in real-time, so as long as you keep the indexer running, it will stay up to date.

If the indexer is stopped for a while, it will need to catch up on the missed blocks. This can take a while, depending on how long it was stopped. Depending on how long it was stopped, you may need an archive node to access the data needed to resync, if this is the case you should consider using the database dump instead.

## Usage

### Requirements

- You will require a private (erigon) archive node to create the database from scratch, alternatively, you can download a database dump from [TODO].
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
[[chains]]
chainID = 1
name = "Ethereum"
shortName = "ETH"
explorerURL = "https://etherscan.io"
routerV2Address = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
factoryV2Address = "0xc0a47dFe034B400B47bDaD5FecDa2621de6c4d95"
routerV3Address = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
factoryV3Address = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
rpcURL = "http://localhost:8545"

[[chains]]
chainID = 56
name = "Binance Smart Chain"
shortName = "BSC"
explorerURL = "https://bscscan.com"
routerV2Address = "0x10ED43C718714eb63d5aA57B78B54704E256024E"
factoryV2Address = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"
routerV3Address = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
factoryV3Address = "0x6725F303b657a9451d8BA641348b6761A6CC7a17"
rpcURL = "http://localhost:8546"

[api]
host = "localhost"
port = 8080
authDefaultExpirary = 7776000 # 90 days
authKeyType = "hex64" # uuid | hex16 | hex32 | hex64 | hex128 | hex256 | jwt
authMasterKey = "my-master-key" # key to access auth methods
authProvider = "sql" # sql | memory | noauth
rateLimitRequests = 500 # max requests per minute
rateLimitStrategy = "ip" # ip | key (requires auth)

[sync.pairs]
batchConcurrency = 2
batchSize = 10
blockRange = 200

[sync.tokens]
batchConcurrency = 2
batchSize = 10

[sync.blockTimestamps]
batchConcurrency = 2
batchSize = 10

# currently only postgres is supported
[storage.postgres]
host = "localhost"
password = "postgres"
port = "5432"
sslmode = "disable"
user = "postgres"

```

### Running

```bash
go run ./cmd/api/main.go --config config.toml
```

### Public API

The API is JSON-RPC 2.0 compliant and is served on port 8080 by default.
The available methods are:

- `idx_getBlockNumber` - Get the current block number for a given chain

- `idx_getChains` - Get the chain IDs for the supported chains

- `idx_getBlockTimestamps` - Get the timestamp for a range of block numbers

- `idx_getBlockAtTimestamp` - Get the block number at a timestamp

- `idx_findTokens` - Find tokens by using find params

- `idx_getTokenCount` - Get the total number of tokens

- `idx_findPairs` - Find pairs by using find params

- `idx_getPairCount` - Get the total number of pairs

- `idx_getWalletBalances` - Get wallet balances for a pair (WIP)

- `idx_getTokenHolders` - Get token holders for a token (WIP)

- `idx_getOHLCVT` - Get OHLCV chart data for a pair (WIP)

### Private API

Private API methods require the Master API key to be set in the config file.
The available methods are:

- `auth_generateKey` - Generate a new API key

- `auth_deleteKey` - Delete an API key

- `auth_getKeyStats` - Get usage information for an API key

- `auth_getAuthMethod` - Get the current auth method

- `auth_getKeyType` - Get the type of API keys used for auth (uuid, hex32, hex64 ...etc)

## JSON-RPC API

### inxi_getBlockNumber

Get the current block numbers for all configured chains

#### parameters: none

#### Example Request

```json
{
  "jsonrpc": "2.0",
  "method": "idx_getBlockNumber",
  "id": "1"
}
```

#### Example Response

```json
{
  "id": "1",
  "method": "idx_getBlockNumber",
  "result": {
    "1": 17461068,
    "56": 32868781
  }
}
```

### idx_getChains

Get the chain info for all configured chains

#### parameters: none

#### Example Request

```json
{
  "jsonrpc": "2.0",
  "method": "idx_getChains",
  "id": "1"
}
```

#### Example Response

```json
{
  "id": "1",
  "method": "idx_getChains",
  "result": [
    {
      "chain_id": 1,
      "name": "Ethereum",
      "short_name": "ETH",
      "explorer_url": "https://etherscan.io"
    },
    {
      "chain_id": 56,
      "name": "Binance Smart Chain",
      "short_name": "BSC",
      "explorer_url": "https://bscscan.com"
    }
  ]
}
```

### idx_getBlockTimestamps

Get the block timestamps for a range of block numbers

#### Parameters:

| Parameter    | Type  | Description                            |
| ------------ | ----- | -------------------------------------- |
| `chain_id`   | int64 | The blockchain network ID.             |
| `from_block` | int64 | The starting block number (inclusive). |
| `to_block`   | int64 | The ending block number (inclusive).   |

`from_block` and `to_block` are inclusive, meaning to_block and from_block will be included in the results. to_block - from_block + 1 will be the number of results.

#### Example Request

```json
{
  "jsonrpc": "2.0",
  "method": "idx_getBlockTimestamps",
  "params": {
    "chain_id": 1,
    "from_block": 17000000,
    "to_block": 17000010
  },
  "id": "1"
}
```

#### Example Response

```json
{
  "id": "1",
  "method": "idx_getBlockTimestamps",
  "result": [
    {
      "block": 17000000,
      "timestamp": 1680911891
    },
    {
      "block": 17000001,
      "timestamp": 1680911903
    },
    {
      "block": 17000002,
      "timestamp": 1680911915
    }
  ]
}
```

### idx_getBlockAtTimestamp

Get the block number closest to a given timestamp

#### Parameters:

| Parameter   | Type  | Description                |
| ----------- | ----- | -------------------------- |
| `chain_id`  | int64 | The blockchain network ID. |
| `timestamp` | int64 | Timestamp                  |

#### Example Request

```json
{
  "jsonrpc": "2.0",
  "method": "idx_getBlockAtTimestamp",
  "params": {
    "chain_id": 1,
    "timestamp": 1680911893
  },
  "id": "1"
}
```

#### Example Response

```json
{
  "id": "1",
  "method": "idx_getBlockAtTimestamp",
  "result": {
    "block": 17000000,
    "timestamp": 1680911891
  }
}
```

### `idx_findTokens`

Find tokens using various filters and options.

#### Parameters:

| Parameter  | Type   | Description                |
| ---------- | ------ | -------------------------- |
| `chain_id` | int64  | The blockchain network ID. |
| `filter`   | Object | Criteria to filter tokens. |
| `options`  | Object | Additional query options.  |

#### `TokenFilter` Object:

| Field        | Type   | Description                                   |
| ------------ | ------ | --------------------------------------------- |
| `address`    | string | The token's address.                          |
| `creator`    | string | The creator's address.                        |
| `from_block` | int64  | The starting block number for token creation. |
| `to_block`   | int64  | The ending block number for token creation.   |
| `name`       | string | The name of the token.                        |
| `symbol`     | string | The token's symbol.                           |
| `decimals`   | uint8  | The number of decimals for the token.         |
| `fuzzy`      | bool   | Enable fuzzy search for string fields.        |

#### `Options` Object:

| Field        | Type   | Description                            |
| ------------ | ------ | -------------------------------------- |
| `limit`      | int64  | Maximum number of results to return.   |
| `offset`     | int64  | Offset for pagination.                 |
| `sort_by`    | string | Field to sort the results by.          |
| `sort_order` | string | Order to sort the results (asc, desc). |

#### Sortable Fields:

- `address`
- `creator`
- `name`
- `symbol`
- `decimals`
- `created_at`

#### Sort Orders:

- `asc`
- `desc`

#### Example Request

```json
{
  "jsonrpc": "2.0",
  "method": "idx_findTokens",
  "params": {
    "chain_id": 1,
    "filter": {
      "name": "I am not in danger",
      "symbol": "walter",
      "fuzzy": true
    },
    "options": {
      "limit": 100,
      "sort_by": "created_at",
      "sort_order": "desc"
    }
  },
  "id": "23"
}
```

#### Example Response

```json
{
  "id": "23",
  "method": "idx_findTokens",
  "result": [
    {
      "address": "0x50E7e4E7fa109A59B255bE882846f4186677f406",
      "name": "Who are you talking to right now? Who is it you think you see? Do you know how much I make a year? I mean, even if I told you, you wouldn't believe it. Do you know what would happen if I suddenly decided to stop going into work? A business big enough that it could be listed on the NASDAQ goes belly up. Disappears. It ceases to exist, without me. No, you clearly don't know who you're talking to, so let me clue you in. I am not in danger, Skyler. I AM the danger. A guy opens his door and gets shot, and you think that of me? No! I am the one who knocks!",
      "symbol": "WALTER",
      "decimals": 9,
      "creator": "0x8e89ac066DE630Db9658aB5FA8FeB4ae85279b30",
      "created_at": 17931545,
      "creation_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "chain_id": 1
    },
    {
      "address": "0x0D8da06819bC5bf57cBDC1C9E499F3B3982584Ac",
      "name": "I am not in danger, Skyler. I am the danger. A guy opens his door and gets shot, and you think that of me? No! I am the one who knocks!",
      "symbol": "WALTER",
      "decimals": 9,
      "creator": "0x254fFf07998de67cF68e9e4CB0dC075430c01eFd",
      "created_at": 15156104,
      "creation_hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "chain_id": 1
    }
  ]
}
```

### `idx_getTokenCount`

Get the total number of tokens for a given chain.

#### Parameters:

| Parameter  | Type  | Description                |
| ---------- | ----- | -------------------------- |
| `chain_id` | int64 | The blockchain network ID. |

#### Example Request

```json
{
  "jsonrpc": "2.0",
  "method": "idx_getTokenCount",
  "params": {
    "chain_id": 56
  },
  "id": "1"
}
```

#### Example Response

```json
{
  "id": "1",
  "method": "idx_getTokenCount",
  "result": 1418709
}
```

### `idx_findPairs`

Find pairs using various filters and options.

#### Parameters:

| Parameter  | Type   | Description                |
| ---------- | ------ | -------------------------- |
| `chain_id` | int64  | The blockchain network ID. |
| `filter`   | Object | Criteria to filter pairs.  |
| `options`  | Object | Additional query options.  |

#### `PairFilter` Object:

| Field            | Type   | Description                                        |
| ---------------- | ------ | -------------------------------------------------- |
| `token0_address` | string | The address of token0                              |
| `token1_address` | string | The address of token1                              |
| `pool_address`   | string | The address of the LP                              |
| `from_block`     | int64  | The starting block                                 |
| `to_block`       | int64  | The ending block                                   |
| `fee`            | int64  | The fee of the pair (v3 only)                      |
| `tick_spacing`   | int64  | The tick spacing of the pair (v3 only)             |
| `hash`           | string | The hash of the pair                               |
| `pool_type`      | uint8  | The pool type of the pair (`2` for v2, `3` for v3) |
| `fuzzy`          | bool   | Enable fuzzy search for string fields.             |

#### `Options` Object:

| Field        | Type   | Description                            |
| ------------ | ------ | -------------------------------------- |
| `limit`      | int64  | Maximum number of results to return.   |
| `offset`     | int64  | Offset for pagination.                 |
| `sort_by`    | string | Field to sort the results by.          |
| `sort_order` | string | Order to sort the results (asc, desc). |

#### Sortable Fields:

- `token0_address`
- `token1_address`
- `pool_address`
- `fee`
- `tick_spacing`
- `hash`
- `pool_type`
- `created_at`

#### Sort Orders:

- `asc`
- `desc`

#### Example Request

```json
{
  "jsonrpc": "2.0",
  "method": "idx_findPairs",
  "params": {
    "chain_id": 1,
    "filter": {
      "token1_address": "0xcf299bd11ceceeed13e0c6d155e70240de11e059"
    },
    "options": {
      "limit": 1,
      "sort_by": "created_at",
      "sort_order": "desc"
    }
  },
  "id": "1"
}
```

#### Example Response

```json
{
  "id": "1",
  "method": "idx_findPairs",
  "result": [
    {
      "token0_address": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
      "token1_address": "0xcf299bd11ceceeed13e0c6d155e70240de11e059",
      "fee": 0,
      "tick_spacing": 0,
      "pool_address": "0xf8a8d7bbc800007b4b9325ac4938b5e0ac24002b",
      "pool_type": 2,
      "created_at": 17991353,
      "hash": "0xf2d398d34ff648c358d792e673d786c2ea0a434d27e8a316d7ba3b792cd7300c",
      "chain_id": 1
    }
  ]
}
```

### `idx_getPairCount`

Get the total number of pairs for a given chain.

#### Parameters:

| Parameter  | Type  | Description                |
| ---------- | ----- | -------------------------- |
| `chain_id` | int64 | The blockchain network ID. |

#### Example Request

```json
{
  "jsonrpc": "2.0",
  "method": "idx_getPairCount",
  "params": {
    "chain_id": 56
  },
  "id": "1"
}
```

#### Example Response

```json
{
  "id": "1",
  "method": "idx_getPairCount",
  "result": 1418709
}
```

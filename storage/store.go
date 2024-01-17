package store

import "github.com/autoapev1/indexer/types"

type Store interface {
	// block timestamp
	GetTimestampAtBlock(int64) (types.BlockTimestamp, error)
	SetBlockTimestamp(types.BlockTimestamp) error
	BulkSetBlockTimestamp([]types.BlockTimestamp) error
	BulkGetBlockTimestamp(to int, from int) ([]types.BlockTimestamp, error)

	// token info
	GetTokenInfo(string) (*types.TokenInfo, error)
	InsertTokenInfo(*types.TokenInfo) error
	BulkInsertTokenInfo([]*types.TokenInfo) error

	// pair info
	GetPairInfoByPair(string) (*types.PairInfo, error)
	GetPairsWithToken(string) ([]*types.PairInfo, error)
	SetPairInfo(*types.PairInfo) error
	BulkInsertPairInfo([]*types.PairInfo) error

	// util
	GetUniqueAddressesFromPairs() ([]string, error)
	GetUniqueAddressesFromTokens() ([]string, error)
}

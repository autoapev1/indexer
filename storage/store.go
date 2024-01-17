package storage

import "github.com/autoapev1/indexer/types"

type Store interface {
	// block timestamp
	GetTimestampAtBlock(int64) (*types.BlockTimestamp, error)
	SetBlockTimestamp(*types.BlockTimestamp) error
	BulkSetBlockTimestamp([]*types.BlockTimestamp) error
	BulkGetBlockTimestamp(to int, from int) ([]*types.BlockTimestamp, error)

	// token info
	GetTokenInfo(string) (*types.Token, error)
	InsertTokenInfo(*types.Token) error
	BulkInsertTokenInfo([]*types.Token) error

	// pair info
	GetPairInfoByPair(string) (*types.Pair, error)
	GetPairsWithToken(string) ([]*types.Pair, error)
	SetPairInfo(*types.Pair) error
	BulkInsertPairInfo([]*types.Pair) error

	// util
	GetUniqueAddressesFromPairs() ([]string, error)
	GetUniqueAddressesFromTokens() ([]string, error)
}

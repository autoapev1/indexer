package storage

import "github.com/autoapev1/indexer/types"

type Store interface {
	GetChainID() int64

	// block timestamp
	GetTimestampAtBlock(int64) (*types.BlockTimestamp, error)
	GetBlockAtTimestamp(int64) (*types.BlockTimestamp, error)
	GetHight() (int64, error)
	InsertBlockTimestamp(*types.BlockTimestamp) error
	BulkInsertBlockTimestamp([]*types.BlockTimestamp) error
	BulkGetBlockTimestamp(to int64, from int64) ([]*types.BlockTimestamp, error)

	// token info
	GetTokenInfo(string) (*types.Token, error)
	GetTokenCount() (int64, error)
	InsertTokenInfo(*types.Token) error
	BulkInsertTokenInfo([]*types.Token) error

	// pair info
	GetPairInfoByPair(string) (*types.Pair, error)
	GetPairsWithToken(string) ([]*types.Pair, error)
	GetPairCount() (int64, error)
	InsertPairInfo(*types.Pair) error
	BulkInsertPairInfo([]*types.Pair) error

	// util
	GetUniqueAddressesFromPairs() ([]string, error)
	GetUniqueAddressesFromTokens() ([]string, error)
}

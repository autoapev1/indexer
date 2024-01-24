package storage

import "github.com/autoapev1/indexer/types"

type Store interface {
	Init() error
	Ready() bool
	GetChainID() int64
	GetHight() (int64, error)

	// block timestamp
	GetBlockAtTimestamp(int64) (*types.BlockTimestamp, error)
	InsertBlockTimestamp(*types.BlockTimestamp) error
	BulkInsertBlockTimestamp([]*types.BlockTimestamp) error
	GetBlockTimestamps(to int64, from int64) ([]*types.BlockTimestamp, error)

	// token info
	FindTokens(*types.FindTokensRequest) ([]*types.Token, error)
	GetTokenCount() (int64, error)
	InsertTokenInfo(*types.Token) error
	BulkInsertTokenInfo([]*types.Token) error

	// pair info
	FindPairs(*types.FindPairsRequest) ([]*types.Pair, error)
	GetPairCount() (int64, error)
	InsertPairInfo(*types.Pair) error
	BulkInsertPairInfo([]*types.Pair) error

	// util
	GetUniqueAddressesFromPairs() ([]string, error)
	GetUniqueAddressesFromTokens() ([]string, error)
	GetPairsWithoutTokenInfo() ([]string, error)
	GetHeights() (*types.Heights, error)
}

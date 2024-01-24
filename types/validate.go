package types

import (
	"errors"
)

var (
	errEmptyRequest       = errors.New("empty request")
	errInvalidChainID     = errors.New("invalid chain_id")
	errMissingChainID     = errors.New("missing required parameter: chain_id")
	errMissingFromBlock   = errors.New("missing required parameter: from_block")
	errMissingToBlock     = errors.New("missing required parameter: to_block")
	errMissingTimestamp   = errors.New("missing required parameter: timestamp")
	errMissingFilter      = errors.New("missing required parameter: filter")
	errInvalidPairSortBy  = errors.New("invalid parameter: sort_by - must be either 'token0_address', 'token1_address', 'pool_address', 'fee', 'tick_spacing', 'hash', 'pool_type', 'created_at'")
	errInvalidTokenSortBy = errors.New("invalid parameter: sort_by - must be either 'address', 'name', 'symbol', 'decimals', 'creator', 'created_at', 'creation_hash'")
	errInvalidSortOrder   = errors.New("invalid parameter: sort_order - must be either 'asc' or 'desc'")
)

func (r *GetBlockTimestampsRequest) Validate() error {
	if r == nil {
		return errEmptyRequest
	}

	if r.ChainID == nil {
		return errMissingChainID
	}

	if *r.ChainID == 0 {
		return errInvalidChainID
	}

	if r.FromBlock == nil {
		return errMissingFromBlock
	}

	if r.ToBlock == nil {
		return errMissingToBlock
	}

	if *r.FromBlock > *r.ToBlock {
		return errors.New("from_block must be less than or equal to to_block")
	}

	if *r.FromBlock < 0 {
		return errors.New("from_block must be greater than or equal to 0")
	}

	if *r.ToBlock < 0 {
		return errors.New("to_block must be greater than or equal to 0")
	}

	if *r.ToBlock == 0 && *r.FromBlock == 0 {
		return errors.New("from_block and to_block cannot both be 0")
	}

	if *r.ToBlock-*r.FromBlock > 10000 {
		return errors.New("from_block and to_block must be within 10000 blocks of each other")
	}

	return nil
}

func (r *GetBlockAtTimestampRequest) Validate() error {
	if r == nil {
		return errEmptyRequest
	}

	if r.ChainID == nil {
		return errMissingChainID
	}

	if *r.ChainID == 0 {
		return errInvalidChainID
	}

	if r.Timestamp == nil {
		return errMissingTimestamp
	}

	if *r.Timestamp <= 0 {
		return errors.New("timestamp must be greater than 0")
	}

	return nil
}

func (r *GetTokenCountRequest) Validate() error {
	if r == nil {
		return errEmptyRequest
	}

	if r.ChainID == nil {
		return errMissingChainID
	}

	if *r.ChainID == 0 {
		return errInvalidChainID
	}

	return nil
}

func (r *GetPairCountRequest) Validate() error {
	if r == nil {
		return errEmptyRequest
	}

	if r.ChainID == nil {
		return errMissingChainID
	}

	if *r.ChainID == 0 {
		return errInvalidChainID
	}

	return nil
}

func (r *FindTokensRequest) Validate() error {
	if r == nil {
		return errEmptyRequest
	}

	if r.ChainID == nil {
		return errMissingChainID
	}

	if *r.ChainID == 0 {
		return errInvalidChainID
	}

	if r.Filter == nil {
		return errMissingFilter
	}

	if r.Filter.ToBlock != nil && r.Filter.FromBlock != nil {
		if *r.Filter.ToBlock < *r.Filter.FromBlock {
			return errors.New("to_block must be greater than or equal to from_block")
		}

		if *r.Filter.ToBlock-*r.Filter.FromBlock > 10000 {
			return errors.New("from_block and to_block must be within 10000 blocks of each other")
		}

		if *r.Filter.FromBlock < 0 {
			return errors.New("from_block must be greater than or equal to 0")
		}

		if *r.Filter.ToBlock < 0 {
			return errors.New("to_block must be greater than or equal to 0")
		}

		if *r.Filter.ToBlock == 0 && *r.Filter.FromBlock == 0 {
			return errors.New("from_block and to_block cannot both be 0")
		}

	}

	if r.Options == nil {
		// use default options
		r.Options = &TokenOptions{
			Offset:    0,
			Limit:     1000,
			SortBy:    TokenSortBy(PairSortByCreatedAt),
			SortOrder: SortASC,
		}
	}

	if r.Options.Offset < 0 {
		return errors.New("offset must be greater than or equal to 0")
	}

	if r.Options.Limit < 0 {
		return errors.New("limit must be greater than or equal to 0")
	}

	if r.Options.Limit == 0 {
		r.Options.Limit = 1000
	}

	if r.Options.Limit > 10000 {
		return errors.New("limit must be less than or equal to 10000")
	}

	if !ValidateSortOrder(r.Options.SortOrder) {
		return errInvalidSortOrder
	}

	if !ValidateTokenSortBy(r.Options.SortBy) {
		return errInvalidTokenSortBy
	}

	return nil
}

func (r *FindPairsRequest) Validate() error {
	if r == nil {
		return errEmptyRequest
	}

	if r.ChainID == nil {
		return errMissingChainID
	}

	if *r.ChainID == 0 {
		return errInvalidChainID
	}

	if r.Filter == nil {
		r.Filter = &PairFilter{}
	}

	if r.Filter.ToBlock != nil && r.Filter.FromBlock != nil {
		if *r.Filter.ToBlock < *r.Filter.FromBlock {
			return errors.New("to_block must be greater than or equal to from_block")
		}

		if *r.Filter.ToBlock-*r.Filter.FromBlock > 10000 {
			return errors.New("from_block and to_block must be within 10000 blocks of each other")
		}

		if *r.Filter.FromBlock < 0 {
			return errors.New("from_block must be greater than or equal to 0")
		}

		if *r.Filter.ToBlock < 0 {
			return errors.New("to_block must be greater than or equal to 0")
		}

		if *r.Filter.ToBlock == 0 && *r.Filter.FromBlock == 0 {
			return errors.New("from_block and to_block cannot both be 0")
		}

	}

	if r.Options == nil {
		// use default options
		r.Options = &PairOptions{
			Offset:    0,
			Limit:     1000,
			SortBy:    PairSortBy(PairSortByCreatedAt),
			SortOrder: SortASC,
		}
	}

	if r.Options.Offset < 0 {
		return errors.New("offset must be greater than or equal to 0")
	}

	if r.Options.Limit < 0 {
		return errors.New("limit must be greater than or equal to 0")
	}

	if r.Options.Limit == 0 {
		r.Options.Limit = 1000
	}

	if r.Options.Limit > 10000 {
		return errors.New("limit must be less than or equal to 10000")
	}

	if !ValidateSortOrder(r.Options.SortOrder) {
		return errInvalidSortOrder
	}

	if !ValidatePairSortBy(r.Options.SortBy) {
		return errInvalidPairSortBy
	}

	return nil
}

func (r *GetHeightsRequest) Validate() error {
	if r == nil {
		return errEmptyRequest
	}

	if r.ChainID == nil {
		return errMissingChainID
	}

	if *r.ChainID == 0 {
		return errInvalidChainID
	}

	return nil
}

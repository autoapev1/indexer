package types

import "errors"

func (r *GetBlockTimestampsRequest) Validate() error {
	if r == nil {
		return errors.New("GetBlockTimestampsRequest is nil")
	}

	if r.ChainID == 0 {
		return errors.New("chain_id is required")
	}

	if r.FromBlock > r.ToBlock {
		return errors.New("from_block must be less than or equal to to_block")
	}

	if r.FromBlock < 0 {
		return errors.New("from_block must be greater than or equal to 0")
	}

	if r.ToBlock < 0 {
		return errors.New("to_block must be greater than or equal to 0")
	}

	if r.ToBlock == 0 && r.FromBlock == 0 {
		return errors.New("from_block and to_block cannot both be 0")
	}

	if r.ToBlock-r.FromBlock > 10000 {
		return errors.New("from_block and to_block must be within 10000 blocks of each other")
	}

	return nil
}

func (r *GetBlockAtTimestampRequest) Validate() error {
	if r == nil {
		return errors.New("GetBlockAtTimestampRequest is nil")
	}

	if r.ChainID == 0 {
		return errors.New("chain_id is required")
	}

	if r.Timestamp <= 0 {
		return errors.New("timestamp must be greater than 0")
	}

	return nil
}

func (r *GetTokenCountRequest) Validate() error {
	if r == nil {
		return errors.New("GetTokenCountRequest is nil")
	}

	if r.ChainID == 0 {
		return errors.New("chain_id is required")
	}

	return nil
}

func (r *GetPairCountRequest) Validate() error {
	if r == nil {
		return errors.New("GetPairCountRequest is nil")
	}

	if r.ChainID == 0 {
		return errors.New("chain_id is required")
	}

	return nil
}

func (r *FindTokensRequest) Validate() error {
	if r == nil {
		return errors.New("FindTokensRequest is nil")
	}

	if r.ChainID == 0 {
		return errors.New("chain_id is required")
	}

	if r.Filter == nil {
		return errors.New("filter is required")
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
		return errors.New("invalid sort_order")
	}

	if !ValidateTokenSortBy(r.Options.SortBy) {
		return errors.New("invalid sort_by")
	}

	return nil
}

func (r *FindPairsRequest) Validate() error {
	if r == nil {
		return errors.New("FindPairsRequest is nil")
	}

	if r.ChainID == 0 {
		return errors.New("chain_id is required")
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
		return errors.New("invalid sort_order")
	}

	if !ValidatePairSortBy(r.Options.SortBy) {
		return errors.New("invalid sort_by")
	}

	return nil
}

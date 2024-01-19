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

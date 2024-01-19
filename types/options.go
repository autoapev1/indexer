package types

type SortOrder string

const (
	SortASC  SortOrder = "asc"
	SortDESC SortOrder = "desc"
)

type TokenSortBy string

const (
	TokenSortByAddress   TokenSortBy = "address"
	TokenSortByCreator   TokenSortBy = "creator"
	TokenSortByName      TokenSortBy = "name"
	TokenSortBySymbol    TokenSortBy = "symbol"
	TokenSortByDecimals  TokenSortBy = "decimals"
	TokenSortByCreatedAt TokenSortBy = "created_at_block"
)

func ValidateSortOrder(order SortOrder) bool {
	switch order {
	case SortASC, SortDESC:
		return true
	default:
		return false
	}
}

func ValidateTokenSortBy(sortBy TokenSortBy) bool {
	switch sortBy {
	case TokenSortByAddress, TokenSortByCreator, TokenSortByName, TokenSortBySymbol, TokenSortByDecimals, TokenSortByCreatedAt:
		return true
	default:
		return false
	}
}

type TokenOptions struct {
	Offset    int64       `json:"offset"`
	Limit     int64       `json:"limit"`
	SortBy    TokenSortBy `json:"sort_by"`
	SortOrder SortOrder   `json:"sort_order"`
}

type PairSortBy string

const (
	PairSortByToken0Address PairSortBy = "token0_address"
	PairSortByToken1Address PairSortBy = "token1_address"
	PairSortByPoolAddress   PairSortBy = "pool_address"
	PairSortByFee           PairSortBy = "fee"
	PairSortByTickSpacing   PairSortBy = "tick_spacing"
	PairSortByHash          PairSortBy = "hash"
	PairSortByPoolType      PairSortBy = "pool_type"
	PairSortByCreatedAt     PairSortBy = "created_at"
)

func ValidatePairSortBy(sortBy PairSortBy) bool {
	switch sortBy {
	case PairSortByToken0Address, PairSortByToken1Address, PairSortByPoolAddress, PairSortByFee, PairSortByTickSpacing, PairSortByHash, PairSortByPoolType, PairSortByCreatedAt:
		return true
	default:
		return false
	}
}

type PairOptions struct {
	Offset    int64      `json:"offset"`
	Limit     int64      `json:"limit"`
	SortBy    PairSortBy `json:"sort_by"`
	SortOrder SortOrder  `json:"sort_order"`
}

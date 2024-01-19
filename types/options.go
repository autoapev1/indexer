package types

type SortOrder string

const (
	SortASC  SortOrder = "asc"
	SortDESC SortOrder = "desc"
)

type TokenSortBy string

const (
	SortByAddress   TokenSortBy = "address"
	SortByCreator   TokenSortBy = "creator"
	SortByName      TokenSortBy = "name"
	SortBySymbol    TokenSortBy = "symbol"
	SortByDecimals  TokenSortBy = "decimals"
	SortByCreatedAt TokenSortBy = "created_at"
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
	case SortByAddress, SortByCreator, SortByName, SortBySymbol, SortByDecimals, SortByCreatedAt:
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
	Token0Address PairSortBy = "token0_address"
	Token1Address PairSortBy = "token1_address"
	PoolAddress   PairSortBy = "pool_address"
	Fee           PairSortBy = "fee"
	TickSpacing   PairSortBy = "tick_spacing"
	Hash          PairSortBy = "hash"
	PoolType      PairSortBy = "pool_type"
	CreatedAt     PairSortBy = "created_at"
)

func ValidatePairSortBy(sortBy PairSortBy) bool {
	switch sortBy {
	case Token0Address, Token1Address, PoolAddress, Fee, TickSpacing, Hash, PoolType, CreatedAt:
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

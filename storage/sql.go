package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/logger"
	"github.com/autoapev1/indexer/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type PostgresStore struct {
	DB      *bun.DB
	ChainID int64
	debug   bool
	ready   bool
}

func NewPostgresDB(conf config.PostgresConfig) *PostgresStore {
	PostgresDB := &PostgresStore{
		ChainID: 0,
	}
	uri := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", conf.User, conf.Password, conf.Host, conf.Name, conf.SSLMode)

	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))

	db := bun.NewDB(pgdb, pgdialect.New())
	PostgresDB.DB = db

	return PostgresDB
}

func (p *PostgresStore) WithChainID(chainID int64) *PostgresStore {
	p.ChainID = chainID
	return p
}

func (p *PostgresStore) WithDebug() *PostgresStore {
	p.debug = true
	return p
}

func (p *PostgresStore) Ready() bool {
	return p.ready
}

func (p *PostgresStore) Init() error {
	st := time.Now()
	err := p.CreateTables()
	if err != nil {
		return err
	}

	p.CreateIndexes()
	if p.debug {
		p.DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	logger.Time("Init()", time.Since(st), true)

	return nil
}

func (p *PostgresStore) CreateTables() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := p.DB.NewCreateTable().
		Model(&types.BlockTimestamp{}).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = p.DB.NewCreateTable().
		IfNotExists().
		Model(&types.Token{}).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = p.DB.NewCreateTable().
		Model(&types.Pair{}).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStore) CreateIndexes() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, _ = p.DB.NewCreateIndex().
		Model(&types.BlockTimestamp{}).
		Column("block").
		Column("timestamp").
		Index("block_timestamp_block_timestamp_idx").
		Exec(ctx)

	_, _ = p.DB.NewCreateIndex().
		Model(&types.Token{}).
		Column("address").
		Index("token_address_idx").
		Exec(ctx)

	_, _ = p.DB.NewCreateIndex().
		Model(&types.Pair{}).
		Column("token0_address").
		Column("token1_address").
		Column("pool_address").
		Index("pair_addresses_idx").
		Exec(ctx)

}

func (p *PostgresStore) GetChainID() int64 {
	return p.ChainID
}

func (p *PostgresStore) GetBlockAtTimestamp(timestamp int64) (*types.BlockTimestamp, error) {
	blockTimestamp := new(types.BlockTimestamp)
	ctx := context.Background()

	const rangeOffset int64 = 20

	err := p.DB.NewSelect().
		Model(blockTimestamp).
		Where("timestamp BETWEEN ? AND ?", timestamp-rangeOffset, timestamp+rangeOffset).
		OrderExpr("ABS(timestamp - ?)", timestamp).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return blockTimestamp, nil
}

func (p *PostgresStore) InsertBlockTimestamp(blockTimestamp *types.BlockTimestamp) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(blockTimestamp).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStore) BulkInsertBlockTimestamp(blockTimestamps []*types.BlockTimestamp) error {
	ctx := context.Background()
	batchSize := 10000

	for i := 0; i < len(blockTimestamps); i += batchSize {
		fmt.Printf("inserting into block_timestamps \tfrom:%d \tto:%d\n", i, i+batchSize)
		end := i + batchSize
		if end > len(blockTimestamps) {
			end = len(blockTimestamps)
		}

		batch := blockTimestamps[i:end]
		_, err := p.DB.NewInsert().Model(&batch).Exec(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			return err
		}
	}
	return nil
}

func (p *PostgresStore) GetBlockTimestamps(to int64, from int64) ([]*types.BlockTimestamp, error) {
	var blockTimestamps []*types.BlockTimestamp
	ctx := context.Background()

	err := p.DB.NewSelect().
		Table("block_timestamps").
		Where("block >= ?", from).
		Where("block <= ?", to).
		Scan(ctx)

	if err != nil {
		return blockTimestamps, err
	}

	return blockTimestamps, nil
}

func (p *PostgresStore) GetHight() (int64, error) {
	var block int64
	ctx := context.Background()
	err := p.DB.NewSelect().
		Table("block_timestamps").
		ColumnExpr("MAX(block)").
		Scan(ctx, &block)
	if err != nil {
		return block, err
	}

	return block, nil
}

func (p *PostgresStore) InsertTokenInfo(tokenInfo *types.Token) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(tokenInfo).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
func (p *PostgresStore) BulkInsertTokenInfo(tokenInfos []*types.Token) error {
	ctx := context.Background()
	batchSize := 10000

	for i := 0; i < len(tokenInfos); i += batchSize {
		tokenInfos[i].Lower()
	}

	for i := 0; i < len(tokenInfos); i += batchSize {
		fmt.Printf("inserting into tokens \tfrom:%d \tto:%d\n", i, i+batchSize)
		end := i + batchSize
		if end > len(tokenInfos) {
			end = len(tokenInfos)
		}

		batch := tokenInfos[i:end]
		_, err := p.DB.NewInsert().Model(&batch).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresStore) GetTokenCount() (int64, error) {
	var count int64
	ctx := context.Background()
	err := p.DB.NewSelect().ColumnExpr("COUNT(*)").Model(&types.Token{}).Scan(ctx, &count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (p *PostgresStore) FindTokens(req *types.FindTokensRequest) ([]*types.Token, error) {
	var tokens []*types.Token

	query := p.DB.NewSelect().Model(&tokens)
	filter := req.Filter

	if filter.Fuzzy {
		if filter.Address != nil && *filter.Address != "" {
			query.Where("address ILIKE ?", fuzWrap(filter.Address))
		}
		if filter.Creator != nil && *filter.Creator != "" {
			query.Where("creator ILIKE ?", fuzWrap(filter.Creator))
		}

		if filter.Name != nil && *filter.Name != "" {
			query.Where("name ILIKE ?", fuzWrap(filter.Name))
		}

		if filter.Symbol != nil && *filter.Symbol != "" {
			query.Where("symbol ILIKE ?", fuzWrap(filter.Symbol))
		}
	} else {
		if filter.Address != nil {
			query.Where("address = ?", filter.Address)
		}
		if filter.Creator != nil {
			query.Where("creator = ?", filter.Creator)
		}

		if filter.Name != nil {
			query.Where("name = ?", filter.Name)
		}

		if filter.Symbol != nil {
			query.Where("symbol = ?", filter.Symbol)
		}
	}

	if filter.Decimals != nil {
		query.Where("decimals = ?", filter.Decimals)
	}

	if filter.FromBlock != nil {
		query.Where("created_at_block >= ?", filter.FromBlock)
	}

	if filter.ToBlock != nil {
		query.Where("created_at_block <= ?", filter.ToBlock)
	}

	if req.Options.SortOrder != "" && req.Options.SortBy != "" {
		query.OrderExpr(fmt.Sprintf("%s %s", req.Options.SortBy, req.Options.SortOrder))
	} else {
		query.OrderExpr("created_at_block ASC")
	}

	if req.Options.Limit != 0 {
		query.Limit(int(req.Options.Limit))
	} else {
		query.Limit(1000)
	}

	if req.Options.Offset != 0 {
		query.Offset(int(req.Options.Offset))
	} else {
		query.Offset(0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := query.Scan(ctx)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

func (p *PostgresStore) InsertPairInfo(pairInfo *types.Pair) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(pairInfo).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStore) BulkInsertPairInfo(pairInfos []*types.Pair) error {
	ctx := context.Background()
	batchSize := 100000

	for i := 0; i < len(pairInfos); i += batchSize {
		pairInfos[i].Lower()
	}

	for i := 0; i < len(pairInfos); i += batchSize {
		fmt.Printf("inserting into pairs \tfrom:%d \tto:%d\n", i, i+batchSize)
		end := i + batchSize
		if end > len(pairInfos) {
			end = len(pairInfos)
		}

		batch := pairInfos[i:end]
		_, err := p.DB.NewInsert().Model(&batch).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresStore) GetPairCount() (int64, error) {
	var count int64
	ctx := context.Background()
	err := p.DB.NewSelect().ColumnExpr("COUNT(*)").Model(&types.Pair{}).Scan(ctx, &count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (p *PostgresStore) FindPairs(req *types.FindPairsRequest) ([]*types.Pair, error) {
	var pairs []*types.Pair

	query := p.DB.NewSelect().Model(&pairs)
	filter := req.Filter

	if filter.Fuzzy {
		if filter.Token0Address != nil && *filter.Token0Address != "" {
			query.Where("token0_address ILIKE ?", fuzWrap(filter.Token0Address))
		}
		if filter.Token1Address != nil && *filter.Token1Address != "" {
			query.Where("token1_address ILIKE ?", fuzWrap(filter.Token1Address))
		}
		if filter.PoolAddress != nil && *filter.PoolAddress != "" {
			query.Where("pool_address ILIKE ?", fuzWrap(filter.PoolAddress))
		}
		if filter.Hash != nil && *filter.Hash != "" {
			query.Where("hash ILIKE ?", fuzWrap(filter.Hash))
		}

	} else {
		if filter.Token0Address != nil {
			query.Where("token0_address = ?", filter.Token0Address)
		}
		if filter.Token1Address != nil {
			query.Where("token1_address = ?", filter.Token1Address)
		}
		if filter.PoolAddress != nil {
			query.Where("pool_address = ?", filter.PoolAddress)
		}
		if filter.Hash != nil {
			query.Where("hash = ?", filter.Hash)
		}
	}

	if filter.FromBlock != nil {
		query.Where("created_at_block >= ?", filter.FromBlock)
	}

	if filter.ToBlock != nil {
		query.Where("created_at_block <= ?", filter.ToBlock)
	}

	if filter.Fee != nil {
		query.Where("fee = ?", filter.Fee)
	}

	if filter.TickSpacing != nil {
		query.Where("tick_spacing = ?", filter.TickSpacing)
	}

	if filter.PoolType != nil {
		query.Where("pool_type = ?", filter.PoolType)
	}

	if req.Options.SortOrder != "" && req.Options.SortBy != "" {
		query.OrderExpr(fmt.Sprintf("%s %s", req.Options.SortBy, req.Options.SortOrder))
	} else {
		query.OrderExpr("created_at_block ASC")
	}

	if req.Options.Limit != 0 {
		query.Limit(int(req.Options.Limit))
	} else {
		query.Limit(1000)
	}

	if req.Options.Offset != 0 {
		query.Offset(int(req.Options.Offset))
	} else {
		query.Offset(0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := query.Scan(ctx)
	if err != nil {
		return pairs, err
	}

	return pairs, nil
}

func (p *PostgresStore) GetUniqueAddressesFromPairs() ([]string, error) {
	// Query to get distinct addresses from both token0 and token1
	var addresses []string
	ctx := context.Background()
	err := p.DB.NewSelect().
		Table("pairs").
		ColumnExpr("DISTINCT token0_address").Scan(ctx, &addresses)
	if err != nil {
		return addresses, err
	}

	err = p.DB.NewSelect().
		Table("pairs").
		ColumnExpr("DISTINCT token1_address").Scan(ctx, &addresses)
	if err != nil {
		return addresses, err
	}

	// Remove duplicates
	uniqueAddresses := make(map[string]bool)
	for _, address := range addresses {
		uniqueAddresses[address] = true
	}

	addresses = []string{}
	for address := range uniqueAddresses {
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (p *PostgresStore) GetUniqueAddressesFromTokens() ([]string, error) {
	var addresses []string
	ctx := context.Background()
	err := p.DB.NewSelect().
		Table("tokens").
		ColumnExpr("DISTINCT address").Scan(ctx, &addresses)
	if err != nil {
		return addresses, err
	}

	return addresses, nil
}

func (p *PostgresStore) GetPairsWithoutTokenInfo() ([]string, error) {
	pairs, err := p.GetUniqueAddressesFromPairs()
	if err != nil {
		return nil, err
	}

	tokenAddresses, err := p.GetUniqueAddressesFromTokens()
	if err != nil {
		return nil, err
	}

	tokens := make(map[string]struct{})
	for _, token := range tokenAddresses {
		tokens[strings.ToLower(token)] = struct{}{}
	}

	var missing []string
	for _, pair := range pairs {
		if _, exists := tokens[strings.ToLower(pair)]; !exists {
			missing = append(missing, pair)
		}
	}

	return missing, nil
}

func (p *PostgresStore) GetHeights() (*types.Heights, error) {
	heights := &types.Heights{
		Blocks: 0,
		Tokens: 0,
		Pairs:  0,
	}

	ctx := context.Background()

	err := p.DB.NewSelect().
		ColumnExpr("MAX(block)").
		Model(&types.BlockTimestamp{}).
		Scan(ctx, &heights.Blocks)
	if err != nil {
		return heights, err
	}

	err = p.DB.NewSelect().
		ColumnExpr("MAX(created_at)").
		Model(&types.Token{}).
		Scan(ctx, &heights.Tokens)
	if err != nil {
		return heights, err
	}

	err = p.DB.NewSelect().
		ColumnExpr("MAX(created_at)").
		Model(&types.Pair{}).
		Scan(ctx, &heights.Pairs)
	if err != nil {
		return heights, err
	}

	return heights, nil
}

var _ Store = &PostgresStore{}

func fuzWrap(s *string) string {
	return fmt.Sprintf("%%%s%%", *s)
}

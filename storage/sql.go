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
}

func NewPostgresDB(conf config.PostgresConfig) *PostgresStore {
	PostgresDB := &PostgresStore{
		ChainID: 0,
	}
	uri := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", conf.User, conf.Password, conf.Host, conf.Name, conf.SSLMode)

	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))

	db := bun.NewDB(pgdb, pgdialect.New())
	PostgresDB.DB = db

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	//db.SetConnMaxIdleTime(30 * time.Second)

	return PostgresDB
}

func (p *PostgresStore) WithChainID(chainID int64) *PostgresStore {
	p.ChainID = chainID
	return p
}

func (p *PostgresStore) Init() error {
	st := time.Now()
	err := p.CreateTables()
	if err != nil {
		return err
	}

	p.CreateIndexes()

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
		Concurrently().
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

func (p *PostgresStore) GetTimestampAtBlock(blockNumber int64) (*types.BlockTimestamp, error) {
	blockTimestamp := new(types.BlockTimestamp)
	ctx := context.Background()

	err := p.DB.NewSelect().Model(blockTimestamp).Where("block = ?", blockNumber).Scan(ctx)
	if err != nil {
		return blockTimestamp, err
	}

	return blockTimestamp, nil
}

func (p *PostgresStore) GetBlockAtTimestamp(timestamp int64) (*types.BlockTimestamp, error) {
	blockTimestamp := new(types.BlockTimestamp)
	ctx := context.Background()

	err := p.DB.NewSelect().
		Model(blockTimestamp).
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
		fmt.Printf("Inserting blocktimestamps %d to %d\n", i, i+batchSize)
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

func (p *PostgresStore) BulkGetBlockTimestamp(to int64, from int64) ([]*types.BlockTimestamp, error) {
	var blockTimestamps []*types.BlockTimestamp
	ctx := context.Background()

	err := p.DB.NewSelect().Model(&blockTimestamps).
		Where("block >= ?", from).
		Where("block <= ?", to).
		Limit(10000).
		Scan(ctx)

	if err != nil {
		return blockTimestamps, err
	}

	return blockTimestamps, nil
}

func (p *PostgresStore) GetHight() (int64, error) {
	var block int64
	ctx := context.Background()
	err := p.DB.NewSelect().Model(&types.BlockTimestamp{}).ColumnExpr("MAX(block)").Scan(ctx, &block)
	if err != nil {
		return block, err
	}

	return block, nil
}

func (p *PostgresStore) GetTokenInfo(address string) (*types.Token, error) {
	tokenInfo := new(types.Token)
	ctx := context.Background()
	err := p.DB.NewSelect().Model(tokenInfo).Where("address = ?", address).Scan(ctx)
	if err != nil {
		return tokenInfo, err
	}

	return tokenInfo, nil
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
		fmt.Printf("Inserting tokens %d to %d\n", i, i+batchSize)
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

func (p *PostgresStore) GetPairInfoByPair(pair string) (*types.Pair, error) {
	pairInfo := new(types.Pair)
	ctx := context.Background()
	err := p.DB.NewSelect().Model(pairInfo).Where("pool_address = ?", pair).Scan(ctx)
	if err != nil {
		return pairInfo, err
	}

	return pairInfo, nil
}

func (p *PostgresStore) GetPairsWithToken(address string) ([]*types.Pair, error) {
	var pairInfos []*types.Pair
	ctx := context.Background()
	err := p.DB.NewSelect().Model(&pairInfos).Where("token0_address = ? OR token1_address = ?", address, address).Scan(ctx)
	if err != nil {
		return pairInfos, err
	}

	return pairInfos, nil
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
		fmt.Printf("Inserting pairs %d to %d\n", i, i+batchSize)
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

func (p *PostgresStore) GetUniqueAddressesFromPairs() ([]string, error) {
	// Query to get distinct addresses from both token0 and token1
	var addresses []string
	ctx := context.Background()
	err := p.DB.NewSelect().ColumnExpr("DISTINCT token0_address").Scan(ctx, &addresses)
	if err != nil {
		return addresses, err
	}

	err = p.DB.NewSelect().ColumnExpr("DISTINCT token1_address").Scan(ctx, &addresses)
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
	err := p.DB.NewSelect().ColumnExpr("DISTINCT address").Scan(ctx, &addresses)
	if err != nil {
		return addresses, err
	}

	return addresses, nil
}

var _ Store = &PostgresStore{}

package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/logger"
	"github.com/autoapev1/indexer/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type PostgresDB struct {
	DB *bun.DB
}

func NewPostgresDB(conf config.Postgres) *PostgresDB {
	PostgresDB := &PostgresDB{}
	uri := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", conf.User, conf.Password, conf.Host, conf.Name, conf.SSLMode)

	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))

	db := bun.NewDB(pgdb, pgdialect.New())
	PostgresDB.DB = db

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return PostgresDB
}

func (p *PostgresDB) Init() error {
	st := time.Now()
	err := p.CreateTables()
	if err != nil {
		return err
	}

	err = p.CreateInexes()
	if err != nil {
		return err
	}

	logger.Time("Init()", time.Since(st), true)

	return nil
}

func (p *PostgresDB) CreateTables() error {
	_, err := p.DB.NewCreateTable().Model(&types.BlockTimestamp{}).IfNotExists().Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = p.DB.NewCreateTable().IfNotExists().Model(&types.TokenInfo{}).Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = p.DB.NewCreateTable().Model(&types.PairInfo{}).IfNotExists().Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) CreateInexes() error {
	_, err := p.DB.NewCreateIndex().
		Model(&types.BlockTimestamp{}).
		Column("block").
		Column("timestamp").
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = p.DB.NewCreateIndex().
		Model(&types.TokenInfo{}).
		Column("address").
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = p.DB.NewCreateIndex().
		Model(&types.PairInfo{}).
		Column("token0").
		Column("token1").
		Column("pool").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) GetTimestampAtBlock(blockNumber int64) (types.BlockTimestamp, error) {
	blockTimestamp := new(types.BlockTimestamp)
	ctx := context.Background()

	err := p.DB.NewSelect().Model(blockTimestamp).Where("block_number = ?", blockNumber).Scan(ctx)
	if err != nil {
		return *blockTimestamp, err
	}

	return *blockTimestamp, nil
}

func (p *PostgresDB) SetBlockTimestamp(blockTimestamp types.BlockTimestamp) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(&blockTimestamp).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) BulkSetBlockTimestamp(blockTimestamps []types.BlockTimestamp) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(&blockTimestamps).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) BulkGetBlockTimestamp(to int, from int) ([]types.BlockTimestamp, error) {
	blockTimestamps := []types.BlockTimestamp{}
	ctx := context.Background()
	err := p.DB.NewSelect().Model(&blockTimestamps).Where("block_number >= ? AND block_number <= ?", from, to).Scan(ctx)
	if err != nil {
		return blockTimestamps, err
	}

	return blockTimestamps, nil
}

func (p *PostgresDB) GetTokenInfo(address string) (*types.TokenInfo, error) {
	tokenInfo := new(types.TokenInfo)
	ctx := context.Background()
	err := p.DB.NewSelect().Model(tokenInfo).Where("address = ?", address).Scan(ctx)
	if err != nil {
		return tokenInfo, err
	}

	return tokenInfo, nil
}

func (p *PostgresDB) InsertTokenInfo(tokenInfo *types.TokenInfo) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(tokenInfo).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) BulkInsertTokenInfo(tokenInfos []*types.TokenInfo) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(&tokenInfos).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) GetPairInfoByPair(pair string) (*types.PairInfo, error) {
	pairInfo := new(types.PairInfo)
	ctx := context.Background()
	err := p.DB.NewSelect().Model(pairInfo).Where("pair = ?", pair).Scan(ctx)
	if err != nil {
		return pairInfo, err
	}

	return pairInfo, nil
}

func (p *PostgresDB) GetPairsWithToken(address string) ([]*types.PairInfo, error) {
	pairInfos := []*types.PairInfo{}
	ctx := context.Background()
	err := p.DB.NewSelect().Model(&pairInfos).Where("token0 = ? OR token1 = ?", address, address).Scan(ctx)
	if err != nil {
		return pairInfos, err
	}

	return pairInfos, nil
}

func (p *PostgresDB) SetPairInfo(pairInfo *types.PairInfo) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(pairInfo).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) BulkInsertPairInfo(pairInfos []*types.PairInfo) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(&pairInfos).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) GetUniqueAddressesFromPairs() ([]string, error) {
	// Query to get distinct addresses from both token0 and token1
	var addresses []string
	ctx := context.Background()
	err := p.DB.NewSelect().ColumnExpr("DISTINCT token0").Scan(ctx, &addresses)
	if err != nil {
		return addresses, err
	}

	err = p.DB.NewSelect().ColumnExpr("DISTINCT token1").Scan(ctx, &addresses)
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

func (p *PostgresDB) GetUniqueAddressesFromTokens() ([]string, error) {
	var addresses []string
	ctx := context.Background()
	err := p.DB.NewSelect().ColumnExpr("DISTINCT address").Scan(ctx, &addresses)
	if err != nil {
		return addresses, err
	}

	return addresses, nil
}

var _ Store = &PostgresDB{}

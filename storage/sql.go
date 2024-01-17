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

	//db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

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
	_, err := p.DB.NewCreateTable().
		Model(&types.BlockTimestamp{}).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = p.DB.NewCreateTable().
		IfNotExists().
		Model(&types.Token{}).
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = p.DB.NewCreateTable().
		Model(&types.Pair{}).
		IfNotExists().
		Exec(context.Background())
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
		Model(&types.Token{}).
		Column("address").
		Exec(context.Background())
	if err != nil {
		return err
	}

	_, err = p.DB.NewCreateIndex().
		Model(&types.Pair{}).
		Column("token0_address").
		Column("token1_address").
		Column("pool_address").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) GetTimestampAtBlock(blockNumber int64) (*types.BlockTimestamp, error) {
	blockTimestamp := new(types.BlockTimestamp)
	ctx := context.Background()

	err := p.DB.NewSelect().Model(blockTimestamp).Where("block_number = ?", blockNumber).Scan(ctx)
	if err != nil {
		return blockTimestamp, err
	}

	return blockTimestamp, nil
}

func (p *PostgresDB) SetBlockTimestamp(blockTimestamp *types.BlockTimestamp) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(blockTimestamp).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) BulkSetBlockTimestamp(blockTimestamps []*types.BlockTimestamp) error {
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

func (p *PostgresDB) BulkGetBlockTimestamp(to int, from int) ([]*types.BlockTimestamp, error) {
	blockTimestamps := []*types.BlockTimestamp{}
	ctx := context.Background()
	err := p.DB.NewSelect().Model(&blockTimestamps).Where("block_number >= ? AND block_number <= ?", from, to).Scan(ctx)
	if err != nil {
		return blockTimestamps, err
	}

	return blockTimestamps, nil
}

func (p *PostgresDB) GetTokenInfo(address string) (*types.Token, error) {
	tokenInfo := new(types.Token)
	ctx := context.Background()
	err := p.DB.NewSelect().Model(tokenInfo).Where("address = ?", address).Scan(ctx)
	if err != nil {
		return tokenInfo, err
	}

	return tokenInfo, nil
}

func (p *PostgresDB) InsertTokenInfo(tokenInfo *types.Token) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(tokenInfo).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
func (p *PostgresDB) BulkInsertTokenInfo(tokenInfos []*types.Token) error {
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

func (p *PostgresDB) GetPairInfoByPair(pair string) (*types.Pair, error) {
	pairInfo := new(types.Pair)
	ctx := context.Background()
	err := p.DB.NewSelect().Model(pairInfo).Where("pair = ?", pair).Scan(ctx)
	if err != nil {
		return pairInfo, err
	}

	return pairInfo, nil
}

func (p *PostgresDB) GetPairsWithToken(address string) ([]*types.Pair, error) {
	pairInfos := []*types.Pair{}
	ctx := context.Background()
	err := p.DB.NewSelect().Model(&pairInfos).Where("token0 = ? OR token1 = ?", address, address).Scan(ctx)
	if err != nil {
		return pairInfos, err
	}

	return pairInfos, nil
}

func (p *PostgresDB) SetPairInfo(pairInfo *types.Pair) error {
	ctx := context.Background()
	_, err := p.DB.NewInsert().Model(pairInfo).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) BulkInsertPairInfo(pairInfos []*types.Pair) error {
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

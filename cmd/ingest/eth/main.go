package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/autoapev1/indexer/adapter"
	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/storage"
)

func main() {
	var (
		configFile string
	)
	flagSet := flag.NewFlagSet("ingest-eth", flag.ExitOnError)
	flagSet.StringVar(&configFile, "config", "config.toml", "")
	flagSet.Parse(os.Args[1:])

	slog.Info("config", "configFile", configFile)
	err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	conf := config.Get().Storage.Postgres
	conf.Name = "ETH"

	db := storage.NewPostgresDB(conf).WithChainID(1)
	err = db.Init()
	if err != nil {
		panic(err)
	}

	var (
		blocks = true
		tokens = true
		pairs  = true
	)

	if blocks {
		err = IngestBlocks(db)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Ingested blocks\n")
		time.Sleep(1 * time.Second)
	}

	if tokens {
		err = IngestTokens(db)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Ingested tokens\n")
		time.Sleep(1 * time.Second)
	}

	if pairs {
		err = IngestPairs(db)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Ingested pairs\n")
		time.Sleep(1 * time.Second)
	}

}

func IngestPairs(pg *storage.PostgresStore) error {
	data, err := adapter.ReadPairs("./adapter/data/eth.pairs.csv")
	if err != nil {
		return err
	}

	for v := range data {
		data[v].ChainID = 1
	}

	err = pg.BulkInsertPairInfo(data)
	if err != nil {
		return err
	}

	return nil
}

func IngestTokens(pg *storage.PostgresStore) error {
	data, err := adapter.ReadTokens("./adapter/data/eth.tokens.csv")
	if err != nil {
		return err
	}

	for v := range data {
		data[v].ChainID = 1
	}

	err = pg.BulkInsertTokenInfo(data)
	if err != nil {
		return err
	}

	return nil
}

func IngestBlocks(pg *storage.PostgresStore) error {
	data, err := adapter.ReadBlockTimestamps("./adapter/data/eth.timestamps.csv")
	if err != nil {
		return err
	}

	err = pg.BulkInsertBlockTimestamp(data)
	if err != nil {
		return err
	}

	return nil
}

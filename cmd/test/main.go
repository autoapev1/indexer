package main

import (
	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/storage"
)

func main() {
	EthConfig := config.Get().Postgres
	EthConfig.Name = "eth"

	BscConfig := config.Get().Postgres
	BscConfig.Name = "bsc"

	edb := storage.NewPostgresDB(EthConfig)

	err := edb.Init()
	if err != nil {
		panic(err)
	}

	bdb := storage.NewPostgresDB(BscConfig)

	err = bdb.Init()
	if err != nil {
		panic(err)
	}

}

package main

import (
	"github.com/autoapev1/indexer/config"
	store "github.com/autoapev1/indexer/storage"
)

func main() {
	db := store.NewPostgresDB(config.Get().Postgres)

	err := db.Init()
	if err != nil {
		panic(err)
	}

}

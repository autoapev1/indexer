package main

import (
	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/storage"
)

func main() {
	EthConfig := config.Get().Storage.Postgres
	EthConfig.Name = "eth"

	BscConfig := config.Get().Storage.Postgres
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

	// type cat struct {
	// 	wiskars int
	// }

	// count := func(cat *cat) int {
	// 	return cat.wiskars
	// }

	// myCat := cat{
	// 	wiskars: 10,
	// }

	// fmt.Println(count(&myCat))

}

package main

import (
	"fmt"

	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/eth"
	"github.com/autoapev1/indexer/types"
)

func main() {
	config.Parse("config.toml")
	// conf := config.Get().Storage.Postgres
	// conf.Name = "ETH"

	// db := storage.NewPostgresDB(conf).WithChainID(1)
	// err := db.Init()
	// if err != nil {
	// 	panic(err)
	// }

	chain := types.Chain{
		Name:      "Ethereum",
		Http:      "http://localhost:7545",
		ShortName: "ETH",
		ChainID:   1,
	}

	eth := eth.NewNetwork(chain, config.Get())

	if err := eth.Init(); err != nil {
		panic(err)
	}

	tokens, err := eth.GetTokenInfo([]string{"0xE0B7927c4aF23765Cb51314A0E0521A9645F0E2A"})
	if err != nil {
		panic(err)
	}

	for _, v := range tokens {
		fmt.Printf("Token %v\n", v)
	}

	// bts, err := eth.GetBlockTimestamps(0, 10)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, v := range bts {
	// 	fmt.Printf("Block %d timestamp %d\n", v.Block, v.Timestamp)
	// }

	// pairAddrs, err := db.GetUniqueAddressesFromPairs()
	// if err != nil {
	// 	panic(err)
	// }

	// tokenAddrs, err := db.GetUniqueAddressesFromTokens()
	// if err != nil {
	// 	panic(err)
	// }

	// uniqueAddrs := make(map[string]struct{})
	// for _, v := range pairAddrs {
	// 	uniqueAddrs[v] = struct{}{}
	// }

	// for _, v := range tokenAddrs {
	// 	uniqueAddrs[v] = struct{}{}
	// }

	// // find addresses in pairs that are not in tokens
	// for _, v := range pairAddrs {
	// 	if _, ok := uniqueAddrs[v]; !ok {
	// 		fmt.Printf("Pair address %s not found in tokens\n", v)
	// 	}
	// }

	// fmt.Printf("Found %d unique addresses\n", len(uniqueAddrs))
	// fmt.Printf("Found %d unique pair addresses\n", len(pairAddrs))
	// fmt.Printf("Found %d unique token addresses\n", len(tokenAddrs))
}

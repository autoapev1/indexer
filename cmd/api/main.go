package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/autoapev1/indexer/api"
	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/storage"
	"github.com/autoapev1/indexer/utils"
)

func main() {
	var (
		configFile string
	)
	flagSet := flag.NewFlagSet("indexer", flag.ExitOnError)
	flagSet.StringVar(&configFile, "config", "config.toml", "")
	flagSet.Parse(os.Args[1:])

	err := config.Parse(configFile)
	if err != nil {
		log.Fatal(err)
	}

	conf := config.Get()
	_ = conf

	storeMap := storage.NewStoreMap()

	for _, v := range conf.Chains {
		c := conf.Storage.Postgres

		if len(v.ShortName) > 25 {
			slog.Warn("ShortName too long", "chainID", v.ChainID, "ShortName", v.ShortName)
			v.ShortName = v.ShortName[:25]
		}

		c.Name = v.ShortName
		db := storage.NewPostgresDB(c)
		storeMap.SetStore(int64(v.ChainID), db)
	}

	server := api.NewServer(conf, storeMap)

	log.Fatal(server.Listen(utils.ToAddress(conf.API.Host, conf.API.Port)))
}

package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/autoapev1/indexer/api"
	"github.com/autoapev1/indexer/auth"
	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/storage"
)

func main() {
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

	authProviderType := auth.ToProvider(conf.API.AuthProvider)
	var authProvider auth.Provider
	switch authProviderType {
	case auth.AuthProviderSql:
		uri := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", conf.User, conf.Password, conf.Host, conf.Name, conf.SSLMode)
		db := auth.NewSqlDB(uri)
		authProvider = auth.NewSqlAuth(db)

	case auth.AuthProviderMemory:
	case auth.AuthProviderNoAuth:
	default:
		slog.Warn("Invalid Auth Provider", "provider", authProviderType)
	}

	server := api.NewServer(conf.Chains, storeMap).WithAuthProvider(authProvider)

	log.Fatal(server.Listen(conf.API.Host))
}

package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

const defaultConfig = `
ethNodeAddr = "http://127.0.0.1:8500"
bscNodeAddr = "http://127.0.0.1:8400"

[tokens]
batchSize 		 = 10
batchConcurrency = 2

[pairs]
batchSize		 = 10
batchConcurrency = 2
blockRange		 = 200

[storage]
driver = "postgres"		# postgres (only supported for now)

[storage.postgres]
user	 = "postgres"
password = "postgres"
host	 = "localhost"
port	 = "5432"
sslmode  = "disable"

[api]
host = "localhost"
port = 8080

authProvider	    = "sql"			    # none / memory / sql
authKeyType		    = "hex64"		    # uuid / hex16 / hex32 / hex64 / hex128 / hex256
authDefaultExpirary = 7776000 		    # 90 days in seconds
authMasterKey 		= "my-master-key"   # used to generate other keys

rateLimitStrategy = "ip" 			    # ip / key / off
rateLimitRequests = 500				    # per second
`

func getDefaultConfig() Config {
	// parse default config
	var config Config
	err := toml.Unmarshal([]byte(defaultConfig), &config)
	if err != nil {
		panic(err)
	}
	return config
}

// Config has global config
var config Config = getDefaultConfig()

type Config struct {
	EthNodeAddr string
	BscNodeAddr string

	Tokens TokensConfig
	Pairs  PairsConfig

	Storage StorageConfig

	Postgres PostgresConfig

	API APIConfig
}

type TokensConfig struct {
	BatchSize        int
	BatchConcurrency int
}

type PairsConfig struct {
	BatchSize        int
	BatchConcurrency int
	BlockRange       int
}

type StorageConfig struct {
	Driver string
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
	Name     string
}

type APIConfig struct {
	Host                string
	Port                int
	AuthProvider        string
	AuthKeyType         string
	AuthDefaultExpirary int64
	AuthMasterKey       string
	RateLimitStrategy   string
	RateLimitRequests   int
}

func Parse(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile("config.toml", []byte(defaultConfig), os.ModePerm); err != nil {
			return err
		}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(b, &config)
	return err
}

func Get() Config {
	return config
}

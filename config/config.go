package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml/v2"
)

const defaultConfig = `
[[chains]]
chainID = 1
name = "Ethereum"
shortName = "ETH"
explorerURL = "https://etherscan.io"
rpcURL = "http://localhost:8545"

[[chains]]
chainID = 56
name = "Binance Smart Chain"
shortName = "BSC"
explorerURL = "https://bscscan.com"
rpcURL = "http://localhost:8546"

[api]
host = "localhost"
port = 8080
authDefaultExpirary = 7776000 # 90 days
authKeyType = "hex64" # uuid | hex16 | hex32 | hex64 | hex128 | hex256 | jwt
authMasterKey = "my-master-key" # key to access auth methods
authProvider = "sql" # sql | memory | noauth
rateLimitRequests = 500 # max requests per minute
rateLimitStrategy = "ip" # ip | key (requires auth)

[sync.pairs]
batchConcurrency = 2
batchSize = 10
blockRange = 200

[sync.tokens]
batchConcurrency = 2
batchSize = 10

[sync.blockTimestamps]
batchConcurrency = 2
batchSize = 10

# currently only postgres is supported
[storage.postgres]
host = "localhost"
password = "postgres"
port = "5432"
sslmode = "disable"
user = "postgres"

`

// func getDefaultConfig() Config {
// 	// parse default config
// 	var config Config
// 	err := toml.Unmarshal([]byte(defaultConfig), &config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Chains: %+v\n", config.Chains)
// 	// fmt.Printf("Sync: %+v\n", config.Sync)
// 	// fmt.Printf("Storage: %+v\n", config.Storage)
// 	// fmt.Printf("API: %+v\n", config.API)
// 	return config
// }

// Config has global config
var config Config

type Config struct {
	Chains []ChainConfig

	Sync SyncConfig

	Storage StorageConfig

	API APIConfig
}

type ChainConfig struct {
	ChainID     int
	Name        string
	ShortName   string
	ExplorerURL string
	RPCURL      string
}

type SyncConfig struct {
	Tokens          TokensSyncConfig
	Pairs           PairsSyncConfig
	BlockTimestamps BlockTimestampsSyncConfig
}

type BlockTimestampsSyncConfig struct {
	BatchSize        int
	BatchConcurrency int
}

type TokensSyncConfig struct {
	BatchSize        int
	BatchConcurrency int
}

type PairsSyncConfig struct {
	BatchSize        int
	BatchConcurrency int
	BlockRange       int
}

type StorageConfig struct {
	Postgres PostgresConfig
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

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
routerV2Address = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
factoryV2Address = "0xc0a47dFe034B400B47bDaD5FecDa2621de6c4d95"
routerV3Address = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
factoryV3Address = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
rpcURL = "http://localhost:8545"
blockDuration = 12

[[chains]]
chainID = 56
name = "Binance Smart Chain"
shortName = "BSC"
explorerURL = "https://bscscan.com"
routerV2Address = "0x10ED43C718714eb63d5aA57B78B54704E256024E"
factoryV2Address = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"
routerV3Address = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
factoryV3Address = "0x6725F303b657a9451d8BA641348b6761A6CC7a17"
rpcURL = "http://localhost:8546"
blockDuration = 3

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
	ChainID          int
	Name             string
	ShortName        string
	ExplorerURL      string
	RouterV2Address  string
	FactoryV2Address string
	RouterV3Address  string
	FactoryV3Address string
	RPCURL           string
	BlockDuration    int
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

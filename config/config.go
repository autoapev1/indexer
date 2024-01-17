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
batchSize        = 10
batchConcurrency = 2

[pairs]
batchSize        = 10
batchConcurrency = 2
blockRange       = 200

[storage]
driver = "postgres"

[postgres]
user     = "postgres"
password = "postgres"
host     = "localhost"
port     = "5432"
sslmode  = "disable"

[api]
host    = "localhost"
port    = 8080
useAuth = true
masterKey  = "my-master-key" # can be used to generate api keys
rateLimitStrategy = "ip" # ip or key
rateLimit = 500 # per second

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

type Token struct {
	BatchSize        int
	BatchConcurrency int
}

type Pairs struct {
	BatchSize        int
	BatchConcurrency int
}
type Storage struct {
	Driver string
}

type Api struct {
	Host              string
	Port              int
	UseAuth           bool   // todo more auth methods
	ApiKey            string // rather than hardcoded
	RateLimitStrategy string
	RateLimit         int
}

type Postgres struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
}

type Config struct {
	EthNodeAddr string
	BscNodeAddr string
	Token       Token
	Pairs       Pairs
	Storage     Storage
	Postgres    Postgres
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

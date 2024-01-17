package types

import (
	"log"
	"os"
	"strings"
	"time"
)

type Chain struct {
	ChainID       int
	ChainName     string
	Symbol        string
	BlockDuration time.Duration
	Http          string
}

const (
	// eth
	UniswapV2    string = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	UniswapV3    string = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	UniFactoryV2 string = "0x5c69bee701ef814a2b6a3edd4b1652cb9cc5aa6f"
	UniFactoryV3 string = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	EthV2Check   string = "0xa4475450b920B706baF5AdB46d97A83F9b538f9D"
	EthV3Check   string = ""

	// binance-smart-chain
	PancakeswapV2        string = "0x10ED43C718714eb63d5aA57B78B54704E256024E"
	PancakeswapV3        string = "0x1b81D678ffb9C0263b24A97847620C99d213eB14"
	PancakeswapFactoryV2 string = "0x5c69bee701ef814a2b6a3edd4b1652cb9cc5aa6f"
	PancakeswapFactoryV3 string = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	BscV2Check           string = "0xd439e0e20f22a4a482ccd93e45af35b0e46faaf2"
	BscV3Check           string = "0xff81b9848845ee11672bb476e3dfab4c379a771e"
)

var (
	// contract ABIs
	UnswapV2ABI         string = readFileToString("./eth/abi/uniswapV2RouterABI.json")
	UnswapV3ABI         string = readFileToString("./eth/abi/uniswapV3RouterABI.json")
	UniswapFactoryV2ABI string = readFileToString("./eth/abi/uniswapV2FactoryABI.json")
	UniswapFactoryV3ABI string = readFileToString("./eth/abi/uniswapV3FactoryABI.json")
	UniswapV2PoolABI    string = readFileToString("./eth/abi/uniswapV2PoolABI.json")
	UniswapV3PoolABI    string = readFileToString("./eth/abi/uniswapV3PoolABI.json")
	PancakeV2ABI        string = readFileToString("./eth/abi/pancakeswapV2RouterABI.json")
	PancakeV3ABI        string = readFileToString("./eth/abi/pancakeswapV3RouterABI.json")
	PancakeFactoryV2ABI string = readFileToString("./eth/abi/pancakeswapV2FactoryABI.json")
	PancakeFactoryV3ABI string = readFileToString("./eth/abi/pancakeswapV3FactoryABI.json")
	PancakeV2PoolABI    string = readFileToString("./eth/abi/pancakeswapV2PoolABI.json")
	PancakeV3PoolABI    string = readFileToString("./eth/abi/pancakeswapV3PoolABI.json")
	Erc20ABI            string = readFileToString("./eth/abi/erc20.json")
	CheckABI            string = readFileToString("./eth/abi/CheckABI.json")
)

func readFileToString(path string) string {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
		return ""

	}

	if strings.Contains(cwd, "eth") {
		path = "../" + path
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return ""
	}

	return string(data)
}

package types

import (
	"log"
	"os"
	"strings"
)

type Chain struct {
	ChainID       int    `json:"chain_id"`
	Name          string `json:"name"`
	ShortName     string `json:"short_name"`
	ExplorerURL   string `json:"explorer_url"`
	RouterV2      string `json:"router_v2"`
	FactoryV2     string `json:"factory_v2"`
	RouterV3      string `json:"router_v3"`
	FactoryV3     string `json:"factory_v3"`
	BlockDuration int64  `json:"block_duration"`
	Http          string `json:"-"`
}

const (
	// eth
	EthV2RouterAddress     string = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	EthV3RouterAddress     string = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	EthV2FactoryAddress    string = "0x5c69bee701ef814a2b6a3edd4b1652cb9cc5aa6f"
	EthV3FactoryAddress    string = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	EthV2TokenCheckAddress string = ""
	EthV3TokenCheckAddress string = ""

	// binance-smart-chain
	BscV2RouterAddress     string = "0x10ED43C718714eb63d5aA57B78B54704E256024E"
	BscV3RouterAddress     string = "0x1b81D678ffb9C0263b24A97847620C99d213eB14"
	BscV2FactoryAddress    string = "0x5c69bee701ef814a2b6a3edd4b1652cb9cc5aa6f"
	BscV3FactoryAddress    string = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	BscV2TokenCheckAddress string = "0xd439e0e20f22a4a482ccd93e45af35b0e46faaf2"
	BscV3TokenCheckAddress string = "0xff81b9848845ee11672bb476e3dfab4c379a771e"
)

var (
	// contract ABIs
	EthV2RouterABI  string = readFileToString("./eth/abi/ETH_V2_Router_ABI.json")
	EthV3RouterABI  string = readFileToString("./eth/abi/ETH_V3_Router_ABI.json")
	EthV2FactoryABI string = readFileToString("./eth/abi/ETH_V2_Factory_ABI.json")
	EthV3FactoryABI string = readFileToString("./eth/abi/ETH_V3_Factory_ABI.json")
	EthV2PoolABI    string = readFileToString("./eth/abi/ETH_V2_Pool_ABI.json")
	EthV3PoolABI    string = readFileToString("./eth/abi/ETH_V3_Pool_ABI.json")
	BsbV2RouterABI  string = readFileToString("./eth/abi/BSC_V2_Router_ABI.json")
	BscV3RouterABI  string = readFileToString("./eth/abi/BSC_V3_Router_ABI.json")
	BscV2FactoryABI string = readFileToString("./eth/abi/BSC_V2_Factory_ABI.json")
	BscV3FactoryABI string = readFileToString("./eth/abi/BSC_V3_Factory_ABI.json")
	BscV2PoolABI    string = readFileToString("./eth/abi/BSC_V2_Pool_ABI.json")
	BscV3PoolABI    string = readFileToString("./eth/abi/BSC_V3_Pool_ABI.json")
	Erc20ABI        string = readFileToString("./eth/abi/ERC20_ABI.json")
	TokenCheckV2ABI string = readFileToString("./eth/abi/TokenCheck_V2_ABI.json")
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

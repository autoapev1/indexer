package eth

import (
	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Network struct {
	Chain  types.Chain
	config config.Config
	Client *ethclient.Client
}

func NewNetwork(c types.Chain, conf config.Config) *Network {
	return &Network{
		Chain:  c,
		config: conf,
	}
}

func (n *Network) Init() error {
	var err error

	n.Client, err = ethclient.Dial(n.Chain.Http)
	if err != nil {
		panic(err)
	}

	return nil
}

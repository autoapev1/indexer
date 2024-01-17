package eth

import (
	"github.com/autoapev1/indexer/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Network struct {
	Chain  types.Chain
	Client *ethclient.Client
}

func NewNetwork(c types.Chain) *Network {
	return &Network{
		Chain: c,
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

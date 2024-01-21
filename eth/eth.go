package eth

import (
	"fmt"

	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
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

func toMethodChecksum(method string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(method))
	sum := hash.Sum(nil)

	return fmt.Sprintf("0x%x", sum[:4])
}

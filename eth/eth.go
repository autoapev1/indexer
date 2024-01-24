package eth

import (
	"fmt"
	"log/slog"

	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/types"
	"github.com/chenzhijie/go-web3"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

type Network struct {
	Chain  types.Chain
	config config.Config
	Client *ethclient.Client
	Web3   *web3.Web3
	ready  bool
}

func NewNetwork(c types.Chain, conf config.Config) *Network {
	return &Network{
		Chain:  c,
		config: conf,
	}
}

func (n *Network) Ready() bool {
	return n.ready
}

func (n *Network) Init() error {
	var err error

	n.Client, err = ethclient.Dial(n.Chain.Http)
	if err != nil {
		slog.Error("Error initilizing eth client", "error", err)
		return err
	}

	w3, err := web3.NewWeb3(n.Chain.Http)
	if err != nil {
		slog.Error("Error initilizing web3 client", "error", err)
		return err
	}

	n.Web3 = w3

	return nil
}

func toMethodChecksum(method string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(method))
	sum := hash.Sum(nil)

	return fmt.Sprintf("0x%x", sum[:4])
}

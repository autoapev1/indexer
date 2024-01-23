package utils

import (
	"github.com/chenzhijie/go-web3/crypto"
	"github.com/ethereum/go-ethereum/common"
)

func TopicToHash(topic string) common.Hash {
	t := []byte(topic)
	sha3Hash := crypto.Keccak256Hash(t)
	commonHash := common.BytesToHash(sha3Hash)
	return commonHash
}

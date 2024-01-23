package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestPairCreated(t *testing.T) {
	topic := "PairCreated(address,address,address,uint256)"
	hash := TopicToHash(topic)
	fmt.Println(hash.String())
	if hash.String() != "0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9" {
		t.Errorf("hash is not correct: %s", hash.String())
	}
}
func TestPoolCreated(t *testing.T) {
	topic := "PoolCreated(address,address,uint24,int24,address)"
	hash := TopicToHash(topic)
	fmt.Println(hash.String())
	if hash.String() != "0x783cca1c0412dd0d695e784568c96da2e9c22ff989357a2e8b1d9b2b4e6b7118" {
		t.Errorf("hash is not correct: %s", hash.String())
	}
}

func TestExtract(t *testing.T) {
	a := "0x00000000000000000000000007da4c5260c678a3acb554bd295b98d313f5502d"
	address := ExtractAddress(a)
	fmt.Println(address.String())
	if strings.ToLower(address.String()) != "0x07da4c5260c678a3acb554bd295b98d313f5502d" {
		t.Errorf("address is not correct: %s", address.String())
	}
}

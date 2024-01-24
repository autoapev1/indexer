package eth

import (
	"context"
	"errors"
	"log/slog"
	"math/big"
	"strings"

	"github.com/autoapev1/indexer/types"
	"github.com/autoapev1/indexer/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type pairMode int

const (
	pairModeV2 pairMode = iota
	pairModeV3
)

func (n *Network) GetPairs(ctx context.Context, to int64, from int64) ([]*types.Pair, error) {
	pairs := make([]*types.Pair, 0)

	bRange := toRange(to, from)
	if err := bRange.validate(); err != nil {
		return nil, err
	}

	batchRange := n.config.Sync.Pairs.BlockRange

	if batchRange > 1000 {
		batchRange = 200
	}

	var (
		V2factoryAddr string
		V3factoryAddr string
		V2factoryABI  string
		V3factoryABI  string
		V2eventSig    common.Hash
		V3eventSig    common.Hash
	)
	switch n.Chain.ChainID {
	case 1:
		V2factoryABI = types.EthV2FactoryABI
		V3factoryABI = types.EthV3FactoryABI
		V2factoryAddr = types.EthV2FactoryAddress
		V3factoryAddr = types.EthV3FactoryAddress

	case 56:
		V2factoryABI = types.BscV2FactoryABI
		V3factoryABI = types.BscV3FactoryABI
		V2factoryAddr = types.BscV2FactoryAddress
		V3factoryAddr = types.BscV3FactoryAddress
	default:
		return nil, errors.New("invalid chain (n.Chain.ChainID)")
	}

	v2factoryDecoder, err := abi.JSON(strings.NewReader(V2factoryABI))
	if err != nil {
		return nil, err
	}

	v3factoryDecoder, err := abi.JSON(strings.NewReader(V3factoryABI))
	if err != nil {
		return nil, err
	}

	V2eventSig = utils.TopicToHash("PairCreated(address,address,address,uint256)")
	V3eventSig = utils.TopicToHash("PoolCreated(address,address,uint24,int24,address)")

	v2s, v2err := n.getPairs(ctx, v2factoryDecoder, V2eventSig, V2factoryAddr, bRange, pairModeV2)
	v3s, v3err := n.getPairs(ctx, v3factoryDecoder, V3eventSig, V3factoryAddr, bRange, pairModeV3)

	if v2err != nil {
		return pairs, v2err
	}

	if v3err != nil {
		return pairs, v3err
	}

	pairs = append(pairs, v2s...)
	pairs = append(pairs, v3s...)

	return pairs, nil
}

func (n *Network) getPairs(ctx context.Context, decoder abi.ABI, signature common.Hash, factory string, bRange blockRange, mode pairMode) ([]*types.Pair, error) {
	var (
		pairs = make([]*types.Pair, 0)
		err   error
	)
	topic := make([][]common.Hash, 0, 1)
	topic = append(topic, []common.Hash{signature})

	filter := ethereum.FilterQuery{
		FromBlock: big.NewInt(bRange.from),
		ToBlock:   big.NewInt(bRange.to),
		Addresses: []common.Address{common.HexToAddress(factory)},
		Topics:    topic,
	}

	logs, err := n.Client.FilterLogs(ctx, filter)
	if err != nil {
		return nil, err
	}

	switch mode {
	case pairModeV2:
		for _, l := range logs {
			if len(l.Topics) != 3 {
				slog.Warn("error decoding v2 PairCreated event", "error", "len(l.Topics) != 3")
				continue
			}

			p := &types.Pair{
				ChainID:       int16(n.Chain.ChainID),
				CreatedAt:     int64(l.BlockNumber),
				Hash:          l.TxHash.String(),
				Token0Address: common.HexToAddress((l.Topics[1].String())).String(),
				Token1Address: common.HexToAddress((l.Topics[2].String())).String(),
				Fee:           0,
				TickSpacing:   0,
				PoolType:      2,
			}

			decoded, err := decoder.Unpack("PairCreated", l.Data)
			if err != nil {
				slog.Warn("error decoding v2 PairCreated event", "error", err)
				continue
			}

			if len(decoded) != 2 {
				slog.Warn("error decoding v2 PairCreated event", "error", "len(decoded) != 2")
				continue
			}

			pair, ok := decoded[0].(common.Address)
			if !ok {
				slog.Warn("error decoding v2 PairCreated event", "error", "pair, ok := decoded[0].(common.Address)")
				continue
			}

			p.PoolAddress = pair.String()
			p.Lower()
			pairs = append(pairs, p)
		}

	case pairModeV3:
		for _, l := range logs {
			if len(l.Topics) != 4 {
				slog.Warn("error decoding v3 PoolCreated event", "error", "len(l.Topics) != 4")
				continue
			}

			p := &types.Pair{
				ChainID:       int16(n.Chain.ChainID),
				CreatedAt:     int64(l.BlockNumber),
				Hash:          l.TxHash.String(),
				Token0Address: common.HexToAddress((l.Topics[1].String())).String(),
				Token1Address: common.HexToAddress((l.Topics[2].String())).String(),
				Fee:           l.Topics[3].Big().Int64(),
				PoolType:      3,
				PoolAddress:   "0x0000000000000000000000000000000000000000",
				TickSpacing:   0,
			}

			decoded, err := decoder.Unpack("PoolCreated", l.Data)
			if err != nil {
				slog.Warn("error decoding v3 PoolCreated event", "error", err)
				continue
			}

			if len(decoded) != 2 {
				slog.Warn("error decoding v3 PoolCreated event", "error", "len(decoded) != 2")
				continue
			}

			tickSpacing, ok := decoded[0].(*big.Int)
			if !ok {
				slog.Warn("error decoding v3 PoolCreated event", "error", "p.TickSpacing, ok = decoded[0].(*big.Int)")
				continue
			}

			poolAddress, ok := decoded[1].(common.Address)
			if !ok {
				slog.Warn("error decoding v3 PoolCreated event", "error", "poolAddress, ok := decoded[1].(common.Address)")
				continue
			}

			p.PoolAddress = poolAddress.String()
			p.TickSpacing = tickSpacing.Int64()

			p.Lower()
			pairs = append(pairs, p)
		}
	default:
		return nil, errors.New("invalid pair mode")
	}

	return pairs, nil
}

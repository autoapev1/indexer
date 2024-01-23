package eth

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/autoapev1/indexer/types"
	"github.com/autoapev1/indexer/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

type blockRange struct {
	to   int64
	from int64
}

func toRange(to int64, from int64) blockRange {
	return blockRange{
		to:   to,
		from: from,
	}
}

func (b blockRange) validate() error {
	if b.to < b.from {
		return errors.New("to must be greater than from")
	}

	return nil
}

func (n *Network) GetTokenInfo(ctx context.Context, tokens []string) ([]*types.Token, error) {
	tokensInfo := make([]*types.Token, 0, len(tokens))

	concurrency := n.config.Sync.Tokens.BatchConcurrency

	if concurrency <= 0 {
		concurrency = 1
	}

	for i := 0; i < len(tokens); i++ {
		tokens[i] = strings.ToLower(tokens[i])
	}

	s1 := n.makeStage1TokenInfoBatches(tokens)

	workers := make(chan int, concurrency)
	var wg sync.WaitGroup
	counter := 0

	for {
		if counter >= len(s1) {
			break
		}

		workers <- 1
		wg.Add(1)
		batch := s1[counter]
		counter++

		go func(batch []rpc.BatchElem) {
			defer func() {
				<-workers
				wg.Done()
			}()

			tis, err := n.getStage1TokenInfoBatch(ctx, batch)
			if err != nil {
				slog.Error("getStage1TokenInfoBatch", "err", err)
				return
			}

			tokensInfo = append(tokensInfo, tis)
		}(batch)
	}
	wg.Wait()

	stage2, batchedTokens := n.makeStage2TokenInfoBatches(tokensInfo, 50)

	wg = sync.WaitGroup{}
	counter = 0
	workers = make(chan int, concurrency)

	for {
		if counter >= len(stage2) {
			break
		}

		workers <- 1
		wg.Add(1)

		batch := stage2[counter]
		batchedTokens := batchedTokens[counter]
		counter++

		go func(batch []rpc.BatchElem, batchedTokens []*types.Token) {
			defer func() {
				<-workers
				wg.Done()
			}()

			err := n.getStage2TokenInfoBatch(ctx, batch, batchedTokens)
			if err != nil {
				slog.Error("getStage2TokenInfoBatch", "err", err)
				return
			}

		}(batch, batchedTokens)
	}

	wg.Wait()

	return tokensInfo, nil
}

func (n *Network) makeStage1TokenInfoBatches(tokens []string) [][]rpc.BatchElem {
	batches := make([][]rpc.BatchElem, 0, len(tokens))

	for i := 0; i < len(tokens); i++ {
		b := make([]rpc.BatchElem, 4)

		// name
		b[0] = rpc.BatchElem{
			Method: "eth_call",
			Args:   []interface{}{map[string]string{"to": tokens[i], "data": toMethodChecksum("name()")}, "latest"},
			Result: new(string),
		}

		// symbol
		b[1] = rpc.BatchElem{
			Method: "eth_call",
			Args:   []interface{}{map[string]string{"to": tokens[i], "data": toMethodChecksum("symbol()")}, "latest"},
			Result: new(string),
		}

		// decimals
		b[2] = rpc.BatchElem{
			Method: "eth_call",
			Args:   []interface{}{map[string]string{"to": tokens[i], "data": toMethodChecksum("decimals()")}, "latest"},
			Result: new(string),
		}

		// creator
		b[3] = rpc.BatchElem{
			Method: "ots_getContractCreator",
			Args:   []interface{}{tokens[i]},
			Result: new(types.Creator),
		}

		batches = append(batches, b)

	}
	return batches
}

func (n *Network) getStage1TokenInfoBatch(ctx context.Context, batch []rpc.BatchElem) (*types.Token, error) {
	token := &types.Token{}

	if err := n.Client.Client().BatchCallContext(ctx, batch); err != nil {
		return nil, err
	}

	for _, b := range batch {
		if b.Error != nil {
			return nil, b.Error
		}
	}

	to := batch[0].Args[0].(map[string]string)["to"]
	token.Address = to

	name, _ := batch[0].Result.(*string)
	symbol, _ := batch[1].Result.(*string)
	decimals, _ := batch[2].Result.(*string)

	creator, ok := batch[3].Result.(*types.Creator)
	if !ok {
		creator = new(types.Creator)
	}

	token.Name = hexToString(name)
	token.Symbol = hexToString(symbol)
	trimmedDecimals := strings.TrimPrefix(*decimals, "0x")
	decodedDecimals, err := strconv.ParseInt(trimmedDecimals, 16, 8)
	if err != nil {
		decodedDecimals = 0
	}
	token.Decimals = uint8(decodedDecimals)

	token.Creator = creator.Creator
	if token.Creator == "" {
		token.Creator = "unknown"
	}
	token.CreationHash = creator.Hash
	if token.CreationHash == "" {
		token.CreationHash = "0x0000000000000000000000000000000000000000000000000000000000000000"
	}

	token.ChainID = int16(n.Chain.ChainID)

	return token, nil

}

func hexToString(hexStr *string) string {
	if hexStr == nil || *hexStr == "" || *hexStr == "0x" {
		return "unknown"
	}
	// Remove the "0x" prefix if it exists
	trimmedHexStr := strings.TrimPrefix(*hexStr, "0x")

	// Decode the hex string
	decodedBytes, err := hex.DecodeString(trimmedHexStr)
	if err != nil {
		return ""
	}

	// Convert byte array to string and trim any null characters
	resultStr := string(decodedBytes)
	resultStr = strings.Trim(resultStr, "\x00")
	return resultStr
}

func (n *Network) makeStage2TokenInfoBatches(tokens []*types.Token, batchSize int) ([][]rpc.BatchElem, [][]*types.Token) {
	batchCount := len(tokens) / batchSize
	if len(tokens)%batchSize != 0 {
		batchCount++
	}

	batches := make([][]rpc.BatchElem, 0, batchCount)
	tokenBatches := make([][]*types.Token, 0, batchCount)
	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}

		batch := make([]rpc.BatchElem, 0, end-i)
		tokenBatch := make([]*types.Token, 0, end-i)
		for j := i; j < end; j++ {
			batch = append(batch, rpc.BatchElem{
				Method: "eth_getTransactionByHash",
				Args:   []interface{}{tokens[j].CreationHash},
				Result: new(types.BlockNumber),
			})
			tokenBatch = append(tokenBatch, tokens[j])

		}

		batches = append(batches, batch)
		tokenBatches = append(tokenBatches, tokenBatch)
	}

	return batches, tokenBatches
}

func (n *Network) getStage2TokenInfoBatch(ctx context.Context, batch []rpc.BatchElem, tokens []*types.Token) error {

	if err := n.Client.Client().BatchCallContext(ctx, batch); err != nil {
		return err
	}

	for _, b := range batch {
		if b.Error != nil {
			return b.Error
		}
	}
	for i, b := range batch {
		blockNumber, ok := b.Result.(*types.BlockNumber)
		if !ok {
			blockNumber = new(types.BlockNumber)
		}
		b := blockNumber.Number
		b = strings.TrimPrefix(b, "0x")
		if b == "" {
			b = "0"
		}
		blockNumberInt, err := strconv.ParseInt(b, 16, 64)
		if err != nil {
			slog.Error("getStage2TokenInfoBatch", "err", err)
			blockNumberInt = 0
		}
		tokens[i].CreatedAtBlock = blockNumberInt
	}

	return nil

}

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
				PoolAddress:   "unknown",
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

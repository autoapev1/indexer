package eth

import (
	"context"
	"encoding/hex"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/autoapev1/indexer/types"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func (n *Network) GetBlockTimestamps(from int64, to int64) ([]*types.BlockTimestamp, error) {
	var blockTimestamps []*types.BlockTimestamp

	batchSize := n.config.Sync.BlockTimestamps.BatchSize
	concurrency := n.config.Sync.BlockTimestamps.BatchConcurrency

	if batchSize <= 0 {
		batchSize = 100
	}

	if concurrency <= 0 {
		concurrency = 2
	}

	batches := n.makeBlockTimestampBatches(from, to, int64(batchSize))

	workers := make(chan int, concurrency)
	var wg sync.WaitGroup
	counter := 0
	for {
		if counter >= len(batches) {
			break
		}

		workers <- 1
		wg.Add(1)
		batch := batches[counter]
		counter++

		go func(batch []rpc.BatchElem) {
			defer func() {
				<-workers
				wg.Done()
			}()

			bts, err := n.getBlockTimestampBatch(batch)
			if err != nil {
				slog.Error("getBlockTimestampBatch", "err", err)
				return
			}

			blockTimestamps = append(blockTimestamps, bts...)
		}(batch)
	}

	wg.Wait()

	return blockTimestamps, nil
}

func (n *Network) makeBlockTimestampBatches(from int64, to int64, batchSize int64) [][]rpc.BatchElem {
	batchCount := (to - from) / batchSize
	if (to-from)%batchSize != 0 {
		batchCount++
	}

	batches := make([][]rpc.BatchElem, 0, batchCount)

	for i := from; i <= to; i += batchSize {
		end := i + batchSize
		if end > to+1 {
			end = to + 1
		}

		batch := make([]rpc.BatchElem, 0, end-i)
		for j := i; j < end; j++ {
			batch = append(batch, rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []interface{}{j, false},
				Result: new(etypes.Header),
			})
		}
		batches = append(batches, batch)
	}

	return batches
}

func (n *Network) getBlockTimestampBatch(batch []rpc.BatchElem) ([]*types.BlockTimestamp, error) {
	var blockTimestamps []*types.BlockTimestamp

	ctx := context.Background()
	if err := n.Client.Client().BatchCallContext(ctx, batch); err != nil {
		return nil, err
	}

	for _, b := range batch {
		if b.Error != nil {
			return nil, b.Error
		}
	}

	for _, b := range batch {
		header := b.Result.(*etypes.Header)
		blockTimestamps = append(blockTimestamps, &types.BlockTimestamp{
			Block:     header.Number.Int64(),
			Timestamp: int64(header.Time),
		})
	}

	return blockTimestamps, nil
}

func (n *Network) GetTokenInfo(tokens []string) ([]*types.Token, error) {
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

			tis, err := n.getStage1TokenInfoBatch(batch)
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

			err := n.getStage2TokenInfoBatch(batch, batchedTokens)
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

func (n *Network) getStage1TokenInfoBatch(batch []rpc.BatchElem) (*types.Token, error) {
	token := &types.Token{}

	ctx := context.Background()
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

func (n *Network) getStage2TokenInfoBatch(batch []rpc.BatchElem, tokens []*types.Token) error {

	ctx := context.Background()
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

package eth

import (
	"context"
	"log/slog"
	"sync"

	"github.com/autoapev1/indexer/types"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func (n *Network) GetBlockTimestamps(ctx context.Context, from int64, to int64) ([]*types.BlockTimestamp, error) {

	blockTimestamps := make([]*types.BlockTimestamp, 0, to-from)

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

			bts, err := n.getBlockTimestampBatch(ctx, batch)
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

func (n *Network) getBlockTimestampBatch(ctx context.Context, batch []rpc.BatchElem) ([]*types.BlockTimestamp, error) {
	var blockTimestamps []*types.BlockTimestamp

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

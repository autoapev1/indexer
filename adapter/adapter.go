package adapter

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"

	"github.com/autoapev1/indexer/types"
)

func ReadPairs(loc string) ([]*types.Pair, error) {
	file, err := os.Open(loc)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var pairs []*types.Pair
	hashMap := make(map[string]struct{})
	for _, record := range records[1:] { // Skipping header
		hash := record[7]
		if _, exists := hashMap[hash]; exists {
			continue // Skip if hash already exists
		}
		hashMap[hash] = struct{}{}

		fee, _ := strconv.ParseInt(record[4], 10, 64)
		tickSpacing, _ := strconv.ParseInt(record[5], 10, 64)
		createdAt, _ := strconv.ParseInt(record[6], 10, 64)
		poolType, _ := strconv.ParseUint(record[2], 10, 8)

		pair := &types.Pair{
			Token0Address: record[0],
			Token1Address: record[1],
			Fee:           fee,
			TickSpacing:   tickSpacing,
			PoolAddress:   record[3],
			PoolType:      uint8(poolType),
			CreatedAt:     createdAt,
			Hash:          hash,
		}
		pairs = append(pairs, pair)
	}

	// sort by created at asc
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].CreatedAt < pairs[j].CreatedAt
	})

	return pairs, nil
}

func ReadBlockTimestamps(loc string) ([]*types.BlockTimestamp, error) {
	file, err := os.Open(loc)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	blockTimestampsMap := make(map[int64]*types.BlockTimestamp)
	for _, record := range records[1:] { // Skipping header
		block, _ := strconv.ParseInt(record[0], 10, 64)
		timestamp, _ := strconv.ParseInt(record[1], 10, 64)

		if _, exists := blockTimestampsMap[block]; !exists {
			blockTimestampsMap[block] = &types.BlockTimestamp{
				Block:     block,
				Timestamp: timestamp,
			}
		}
	}

	var blockTimestamps []*types.BlockTimestamp
	for _, bt := range blockTimestampsMap {
		blockTimestamps = append(blockTimestamps, bt)
	}

	// sort by block asc
	sort.Slice(blockTimestamps, func(i, j int) bool {
		return blockTimestamps[i].Block < blockTimestamps[j].Block
	})

	return blockTimestamps, nil
}

func ReadTokens(loc string) ([]*types.Token, error) {
	file, err := os.Open(loc)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	tokenMap := make(map[string]*types.Token)
	for _, record := range records[1:] { // Skipping header
		decimals, _ := strconv.ParseUint(record[3], 10, 8)
		createdAtBlock, _ := strconv.ParseInt(record[5], 10, 64)

		address := record[0]
		if _, exists := tokenMap[address]; !exists {
			tokenMap[address] = &types.Token{
				Address:        address,
				Name:           record[1],
				Symbol:         record[2],
				Decimals:       uint8(decimals),
				Creator:        record[4],
				CreatedAtBlock: createdAtBlock,
				ChainID:        0, // Set ChainID to 0 for now
			}
		}
	}

	var tokens []*types.Token
	for _, token := range tokenMap {
		tokens = append(tokens, token)
	}

	// sort by created at block asc
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].CreatedAtBlock < tokens[j].CreatedAtBlock
	})

	return tokens, nil
}

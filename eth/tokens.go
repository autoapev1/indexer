package eth

import (
	"context"
	"time"

	"github.com/autoapev1/indexer/logger"
	"github.com/autoapev1/indexer/types"
	"github.com/ethereum/go-ethereum/common"
)

func (n *Network) GetTokenInfo(ctx context.Context, address common.Address) *types.Token {
	var (
		st     = time.Now()
		result = &types.Token{}
		err    error
	)
	defer func() {
		success := err == nil
		logger.Time("GetTokenInfo()", time.Since(st), success)
	}()

	return result
}

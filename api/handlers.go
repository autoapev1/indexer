package api

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/autoapev1/indexer/auth"
	"github.com/autoapev1/indexer/types"
)

func (s *Server) handleJrpcRequest(r *JRPCRequest, authlvl auth.AuthLevel) Response {

	methodPrefix := getMethodPrefix(r.Method)

	if methodPrefix == MethodInvalid {
		return &JRPCResponse{
			ID:      r.ID,
			JSONRPC: "2.0",
			Error: &JRPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}

	if !hasAccess(methodPrefix, authlvl) {
		return &JRPCResponse{
			ID:      r.ID,
			JSONRPC: "2.0",
			Error: &JRPCError{
				Code:    -32800,
				Message: "Unauthorized",
			},
		}
	}

	switch r.Method {
	// global
	case "idx_getBlockNumber":
		return s.getBlockNumber(r)
	case "idx_getChains":
		return notImplemented(r)

	// block timestamps
	case "idx_getBlockTimestamps":
		return notImplemented(r)
	case "idx_getBlockAtTimestamp":
		return notImplemented(r)

	// tokens
	case "idx_findTokens":
		return notImplemented(r)
	case "idx_getTokenCount":
		return notImplemented(r)

	// pairs
	case "idx_findPairs":
		return notImplemented(r)
	case "idx_getPairCount":
		return notImplemented(r)

	// holdings
	case "idx_getWalletBalances":
		return notImplemented(r)
	case "idx_getTokenHolders":
		return notImplemented(r)

	// charts
	case "idx_getOHLCVChartData":
		return notImplemented(r)

	case "auth_generateKey":
		return notImplemented(r)
	case "auth_deleteKey":
		return notImplemented(r)
	case "auth_getKeyStats":
		return notImplemented(r)
	case "auth_getAuthMethod":
		return notImplemented(r)
	case "auth_getKeyType":
		return notImplemented(r)

	default:
		return &JRPCResponse{
			ID:      r.ID,
			JSONRPC: "2.0",
			Error: &JRPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func notImplemented(r *JRPCRequest) *JRPCResponse {
	return &JRPCResponse{
		ID:      r.ID,
		JSONRPC: "2.0",
		Error: &JRPCError{
			Code:    -32701,
			Message: "Method not Implemented",
		},
	}
}

func (s *Server) getBlockNumber(r *JRPCRequest) *types.GetBlockNumberResponse {

	stores := s.stores.GetAll()
	if stores == nil {
		if s.debug {
			slog.Error("failed to get stores")
		}
		return &types.GetBlockNumberResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: "internal server error",
			},
		}
	}

	blockNumbers := map[int64]int64{}
	for _, store := range stores {
		block, err := store.GetHight()
		if err != nil {
			if s.debug {
				slog.Error("failed to get block number", "err", err)
			}
			return &types.GetBlockNumberResponse{
				ID:     r.ID,
				Method: r.Method,
				Error: &types.JRPCError{
					Code:    -32602,
					Message: "internal server error",
				},
			}
		}
		blockNumbers[store.GetChainID()] = block
	}

	return &types.GetBlockNumberResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: nil,
	}
}

func (s *Server) getChains(r *JRPCRequest) *types.GetChainsResponse {
	req := &types.GetChainsRequest{}

	if err := json.Unmarshal(r.Params, req); err != nil {
		return &types.GetChainsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}
	chains := []types.Chain{}
	for _, c := range s.config.Chains {
		tc := types.Chain{
			ChainID:       c.ChainID,
			Name:          c.Name,
			ShortName:     c.ShortName,
			ExplorerURL:   c.ExplorerURL,
			RouterV2:      c.RouterV2Address,
			FactoryV2:     c.FactoryV2Address,
			RouterV3:      c.RouterV3Address,
			FactoryV3:     c.FactoryV3Address,
			BlockDuration: time.Duration(c.BlockDuration) * time.Second,
		}
		chains = append(chains, tc)
	}

	return &types.GetChainsResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: chains,
	}
}

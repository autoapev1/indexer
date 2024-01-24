package api

import (
	"encoding/json"
	"log/slog"

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
		return s.getChains(r)
	case "idx_getHeights":
		return s.getHeights(r)

	// block timestamps
	case "idx_getBlockTimestamps":
		return s.getBlockTimestamps(r)
	case "idx_getBlockAtTimestamp":
		return s.getBlockAtTimestamp(r)

	// tokens
	case "idx_findTokens":
		return s.findTokens(r)
	case "idx_getTokenCount":
		return s.getTokenCount(r)

	// pairs
	case "idx_findPairs":
		return s.findPairs(r)
	case "idx_getPairCount":
		return s.getPairCount(r)

	// holdings
	case "idx_getWalletBalances":
		return notImplemented(r)
	case "idx_getTokenHolders":
		return notImplemented(r)

	// charts
	case "idx_getOHLCVT":
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
		Result: blockNumbers,
	}
}

func (s *Server) getChains(r *JRPCRequest) *types.GetChainsResponse {
	chains := []types.Chain{}
	for _, c := range s.config.Chains {
		tc := types.Chain{
			ChainID:     c.ChainID,
			Name:        c.Name,
			ShortName:   c.ShortName,
			ExplorerURL: c.ExplorerURL,
		}
		chains = append(chains, tc)
	}

	return &types.GetChainsResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: chains,
	}
}

func (s *Server) getHeights(r *JRPCRequest) *types.GetHeightsResponse {
	req := &types.GetHeightsRequest{}

	if r.Params == nil {
		return &types.GetHeightsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errMissingParams.Error(),
			},
		}
	}

	err := json.Unmarshal(r.Params, req)
	if err != nil {
		return &types.GetHeightsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errUnmarshalParams.Error(),
			},
		}
	}

	err = req.Validate()
	if err != nil {
		return &types.GetHeightsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: err.Error(),
			},
		}
	}

	store := s.stores.GetStore(*req.ChainID)
	if store == nil {
		return &types.GetHeightsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: "invalid chain_id",
			},
		}
	}

	hs, err := store.GetHeights()
	if err != nil {
		if s.debug {
			slog.Error("failed to get heights", "err", err)
		}
		return &types.GetHeightsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errInternalServer.Error(),
			},
		}
	}

	return &types.GetHeightsResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: hs,
	}

}

func (s *Server) getBlockTimestamps(r *JRPCRequest) *types.GetBlockTimestampsResponse {
	req := &types.GetBlockTimestampsRequest{}

	if r.Params == nil {
		return &types.GetBlockTimestampsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errMissingParams.Error(),
			},
		}
	}

	err := json.Unmarshal(r.Params, req)
	if err != nil {
		if s.debug {
			slog.Error("failed to unmarshal params", "err", err)
		}
		return &types.GetBlockTimestampsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errUnmarshalParams.Error(),
			},
		}
	}

	err = req.Validate()
	if err != nil {
		return &types.GetBlockTimestampsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: err.Error(),
			},
		}
	}

	store := s.stores.GetStore(*req.ChainID)
	if store == nil {
		return &types.GetBlockTimestampsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: "invalid chain_id",
			},
		}
	}

	blockTimestamps, err := store.GetBlockTimestamps(*req.ToBlock, *req.FromBlock)
	if err != nil {
		if s.debug {
			slog.Error("failed to get block timestamps", "err", err)
		}
		return &types.GetBlockTimestampsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errInternalServer.Error(),
			},
		}
	}

	return &types.GetBlockTimestampsResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: blockTimestamps,
	}
}

func (s *Server) getBlockAtTimestamp(r *JRPCRequest) *types.GetBlockAtTimestampResponse {
	req := &types.GetBlockAtTimestampRequest{}

	if r.Params == nil {
		return &types.GetBlockAtTimestampResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errMissingParams.Error(),
			},
		}
	}

	err := json.Unmarshal(r.Params, req)
	if err != nil {
		return &types.GetBlockAtTimestampResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errUnmarshalParams.Error(),
			},
		}
	}

	err = req.Validate()
	if err != nil {
		return &types.GetBlockAtTimestampResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: err.Error(),
			},
		}
	}

	store := s.stores.GetStore(*req.ChainID)
	if store == nil {
		return &types.GetBlockAtTimestampResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: "invalid chain_id",
			},
		}
	}

	block, err := store.GetBlockAtTimestamp(*req.Timestamp)
	if err != nil {
		if s.debug {
			slog.Error("failed to get block at timestamp", "err", err)
		}
		return &types.GetBlockAtTimestampResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errInternalServer.Error(),
			},
		}
	}

	return &types.GetBlockAtTimestampResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: block,
	}
}
func (s *Server) findTokens(r *JRPCRequest) *types.FindTokensResponse {
	var req *types.FindTokensRequest

	if r.Params == nil {
		return &types.FindTokensResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errMissingParams.Error(),
			},
		}
	}

	err := json.Unmarshal(r.Params, &req)
	if err != nil {
		return &types.FindTokensResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errUnmarshalParams.Error(),
			},
		}
	}

	err = req.Validate()
	if err != nil {
		return &types.FindTokensResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: err.Error(),
			},
		}
	}

	store := s.stores.GetStore(*req.ChainID)
	if store == nil {
		return &types.FindTokensResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: "invalid chain_id",
			},
		}
	}

	tokens, err := store.FindTokens(req)
	if err != nil {
		if s.debug {
			slog.Error("failed to find tokens", "err", err)
		}

		return &types.FindTokensResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errInternalServer.Error(),
			},
		}
	}

	return &types.FindTokensResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: tokens,
	}
}

func (s *Server) getTokenCount(r *JRPCRequest) *types.GetTokenCountResponse {
	req := &types.GetTokenCountRequest{}

	if r.Params == nil {
		return &types.GetTokenCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errMissingParams.Error(),
			},
		}
	}

	err := json.Unmarshal(r.Params, req)
	if err != nil {
		return &types.GetTokenCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errUnmarshalParams.Error(),
			},
		}
	}

	err = req.Validate()
	if err != nil {
		return &types.GetTokenCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: err.Error(),
			},
		}
	}

	store := s.stores.GetStore(*req.ChainID)
	if store == nil {
		return &types.GetTokenCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: "invalid chain_id",
			},
		}
	}

	count, err := store.GetTokenCount()
	if err != nil {
		if s.debug {
			slog.Error("failed to get token count", "err", err)
		}
		return &types.GetTokenCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errInternalServer.Error(),
			},
		}
	}

	return &types.GetTokenCountResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: count,
	}
}

func (s *Server) findPairs(r *JRPCRequest) *types.FindPairsResponse {
	req := &types.FindPairsRequest{}

	if r.Params == nil {
		return &types.FindPairsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errMissingParams.Error(),
			},
		}
	}

	err := json.Unmarshal(r.Params, req)
	if err != nil {
		return &types.FindPairsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errUnmarshalParams.Error(),
			},
		}
	}

	err = req.Validate()
	if err != nil {
		return &types.FindPairsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: err.Error(),
			},
		}
	}

	store := s.stores.GetStore(*req.ChainID)
	if store == nil {
		return &types.FindPairsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: "invalid chain_id",
			},
		}
	}

	pairs, err := store.FindPairs(req)
	if err != nil {
		if s.debug {
			slog.Error("failed to find pairs", "err", err)
		}

		return &types.FindPairsResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errInternalServer.Error(),
			},
		}
	}

	return &types.FindPairsResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: pairs,
	}
}

func (s *Server) getPairCount(r *JRPCRequest) *types.GetPairCountResponse {
	req := &types.GetPairCountRequest{}

	if r.Params == nil {
		return &types.GetPairCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errMissingParams.Error(),
			},
		}
	}

	err := json.Unmarshal(r.Params, req)
	if err != nil {
		return &types.GetPairCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errUnmarshalParams.Error(),
			},
		}
	}

	err = req.Validate()
	if err != nil {
		return &types.GetPairCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: err.Error(),
			},
		}
	}

	store := s.stores.GetStore(*req.ChainID)
	if store == nil {
		return &types.GetPairCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: "invalid chain_id",
			},
		}
	}

	count, err := store.GetPairCount()
	if err != nil {
		if s.debug {
			slog.Error("failed to get pair count", "err", err)
		}
		return &types.GetPairCountResponse{
			ID:     r.ID,
			Method: r.Method,
			Error: &types.JRPCError{
				Code:    -32602,
				Message: errInternalServer.Error(),
			},
		}
	}

	return &types.GetPairCountResponse{
		ID:     r.ID,
		Method: r.Method,
		Result: count,
	}
}

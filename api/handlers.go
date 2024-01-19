package api

import "github.com/autoapev1/indexer/auth"

func (s *Server) handleJrpcRequest(r *JRPCRequest, authlvl auth.AuthLevel) *JRPCResponse {

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
		return notImplemented(r)
	case "idx_getChains":
		return notImplemented(r)

	// block timestamps
	case "idx_getBlockTimestamps":
		return notImplemented(r)
	case "idx_getBlockAtTimestamp":
		return notImplemented(r)

	// tokens
	case "idx_getTokenByAddress":
		return notImplemented(r)
	case "idx_getTokensByCreator":
		return notImplemented(r)
	case "idx_getTokensInBlock":
		return notImplemented(r)
	case "idx_findTokens":
		return notImplemented(r)

	// pairs
	case "idx_getPairByAddress":
		return notImplemented(r)
	case "idx_getaPairsByToken":
		return notImplemented(r)
	case "idx_getPairsInBlock":
		return notImplemented(r)
	case "idx_findPairs":
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

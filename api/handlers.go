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
	case "idx_getChains":

	// block timestamps
	case "idx_getBlockTimestamps":
	case "idx_getBlockAtTimestamp":

	// tokens
	case "idx_getTokenByAddress":
	case "idx_getTokensByCreator":
	case "idx_getTokensInBlock":
	case "idx_findTokens":

	// pairs
	case "idx_getPairByAddress":
	case "idx_getaPairsByToken":
	case "idx_getPairsInBlock":
	case "idx_findPairs":

	// holdings
	case "idx_getWalletBalances":
	case "idx_getTokenHolders":

	// charts
	case "idx_getOHLCVChartData":

	case "auth_generateKey":
	case "auth_deleteKey":
	case "auth_getKeyStats":
	case "auth_getAuthMethod":
	case "auth_getKeyType":

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
	return nil
}

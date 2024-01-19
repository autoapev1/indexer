package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/autoapev1/indexer/auth"
)

func AuthMiddleware(a auth.Provider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if a == nil {
				slog.Warn("Auth Provider is Nil")
				return
			}

			level, err := a.Authenticate(r)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, &JRPCResponse{
					JSONRPC: "2.0",
					Error: &JRPCError{
						Code:    -32800,
						Message: "Unauthorized",
					},
				})
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, auth.AuthKey, level)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

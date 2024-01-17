package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	errUnauthorized   = errors.New("unauthorized")
	errInternalServer = errors.New("internal server error")
)

type apiHandler func(w http.ResponseWriter, r *http.Request) error

func makeAPIHandler(h apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic in api handler", "error", rec)
				writeError(w, errInternalServer)
			}
		}()

		if err := h(w, r); err != nil {
			slog.Error("api handler error", "err", err)
			writeError(w, errInternalServer)
		}
	}
}

func writeError(w http.ResponseWriter, err error) error {
	return writeJSON(w, http.StatusInternalServerError, &JRPCResponse{
		Error: &JRPCError{
			Code:    -32603,
			Message: err.Error(),
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	errDecodeRequestBody = errors.New("could not decode the request body")
	errUnauthorized      = errors.New("unauthorized")
)

type errorResponse struct {
	Error string `json:"error"`
}

func ErrorResponse(err error) errorResponse {
	return errorResponse{
		Error: err.Error(),
	}
}

type apiHandler func(w http.ResponseWriter, r *http.Request) error

func makeAPIHandler(h apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic in api handler", "error", rec)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		if err := h(w, r); err != nil {
			slog.Error("api handler error", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

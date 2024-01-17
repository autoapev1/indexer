package api

import (
	"encoding/json"
	"net/http"

	"github.com/autoapev1/indexer/auth"
	"github.com/autoapev1/indexer/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router   *chi.Mux
	ethStore storage.Store
	bscStore storage.Store
	Auth     auth.Provider
}

// NewServer returns a new server given a Store interface.
func NewServer(estore storage.Store, bstore storage.Store) *Server {
	return &Server{
		ethStore: estore,
		bscStore: bstore,
	}
}

// Listen starts listening on the given address.
func (s *Server) Listen(addr string) error {
	s.initRouter()
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) initRouter() {
	s.router = chi.NewRouter()

	// middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(s.withAPIToken)
	s.router.Use(middleware.RealIP)

	// routes
	s.router.Get("/rpc", makeAPIHandler(s.handleRPC))
	s.router.Get("/status", handleStatus)

}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	status := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(status)
}

// main request handler func
func (s *Server) handleRPC(w http.ResponseWriter, r *http.Request) error {
	var req []*JRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, &JRPCResponse{
			Error: &JRPCError{
				Code:    -32600,
				Message: err.Error(),
			},
		})
	}

	var resp []*JRPCResponse

	// range over the requests and handle them
	for _, r := range req {
		switch r.Method {
		default:
			resp = append(resp, &JRPCResponse{
				ID:      r.ID,
				JSONRPC: "2.0",
				Error: &JRPCError{
					Code:    -32601,
					Message: "Method not found",
				},
			})
		}
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) withAPIToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader != "a2ecf22c-8ec5-4011-9f09-eed4a7bd86e9" {
			writeJSON(w, http.StatusUnauthorized, &JRPCResponse{
				Error: &JRPCError{
					Code:    -32600,
					Message: errUnauthorized.Error(),
				},
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

package api

import (
	"encoding/json"
	"net/http"

	"github.com/autoapev1/indexer/auth"
	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	chains []config.ChainConfig
	stores *storage.StoreMap
	auth   auth.Provider
}

// NewServer returns a new server given a Store interface.
func NewServer(chains []config.ChainConfig, stores *storage.StoreMap) *Server {
	return &Server{
		chains: chains,
		stores: stores,
	}
}

func (s *Server) WithAuthProvider(ap auth.Provider) *Server {
	s.auth = ap
	return s
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
	s.router.Use(middleware.RealIP)

	// routes
	s.router.Get("/rpc", makeAPIHandler(s.handleBase))
	s.router.Get("/auth", makeAPIHandler(s.handleAuth))
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
func (s *Server) handleBase(w http.ResponseWriter, r *http.Request) error {
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
		response := s.handleBaseRequest(r)
		resp = append(resp, response)

	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleAuth(w http.ResponseWriter, r *http.Request) error {
	var req []*JRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, &JRPCResponse{
			Error: &JRPCError{
				Code:    -36000,
				Message: err.Error(),
			},
		})
	}

	var resp []*JRPCResponse
	// range over the requests and handle them
	for _, r := range req {
		response := s.handleAuthRequest(r)
		resp = append(resp, response)

	}

	return writeJSON(w, http.StatusOK, resp)
}

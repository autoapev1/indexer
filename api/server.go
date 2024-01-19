package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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

func (s *Server) initAuthProvider() error {
	conf := config.Get()
	authProviderType := auth.ToProvider(conf.API.AuthProvider)
	var authProvider auth.Provider
	switch authProviderType {

	case auth.AuthProviderSql:
		uri := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s",
			conf.Storage.Postgres.User,
			conf.Storage.Postgres.Password,
			conf.Storage.Postgres.Host,
			conf.Storage.Postgres.Name,
			conf.Storage.Postgres.SSLMode)
		db := auth.NewSqlDB(uri)
		authProvider = auth.NewSqlAuthProvider(db)

	case auth.AuthProviderMemory:

		authProvider = auth.NewMemoryProvider()

	case auth.AuthProviderNoAuth:

		authProvider = auth.NewNoAuthProvider()

	default:
		slog.Warn("Invalid Auth Provider", "provider", authProviderType)
	}

	s.auth = authProvider
	return nil
}

// Listen starts listening on the given address.
func (s *Server) Listen(addr string) error {

	if err := s.initAuthProvider(); err != nil {
		return err
	}

	s.initRouter()

	fmt.Printf("API Server Listening on: \t%s\n", addr)
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) initRouter() {
	s.router = chi.NewRouter()

	// middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(authMiddleware(s.auth))
	s.router.Use(middleware.RealIP)

	// routes
	s.router.Get("/", makeAPIHandler(s.handleRequest))
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
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) error {
	var req []*JRPCRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return writeError(w, http.StatusBadRequest, errReadingBody)
	}

	if len(body) == 0 {
		// Return a blank HTML page
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body></body></html>"))
		return nil
	}

	if err := json.Unmarshal(body, &req); err != nil {
		return writeError(w, http.StatusBadRequest, errUnmarshalRequest)
	}

	authLevel := r.Context().Value(auth.AuthKey).(auth.AuthLevel)
	if !auth.IsValidAuthLevel(authLevel) {
		return writeJSON(w, http.StatusInternalServerError, &JRPCResponse{
			Error: &JRPCError{
				Code:    -32600,
				Message: "internal server error",
			},
		})
	}

	var resp []*JRPCResponse
	// range over the requests and handle them
	for _, r := range req {
		response := s.handleJrpcRequest(r, authLevel)
		resp = append(resp, response)
	}

	return writeJSON(w, http.StatusOK, resp)
}

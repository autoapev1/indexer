package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/autoapev1/indexer/auth"
	"github.com/autoapev1/indexer/config"
	"github.com/autoapev1/indexer/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router    *chi.Mux
	config    config.Config
	stores    *storage.StoreMap
	auth      auth.Provider
	rateLimit *rateLimiter
	debug     bool
}

// NewServer returns a new server given a Store interface.
func NewServer(conf config.Config, stores *storage.StoreMap) *Server {
	return &Server{
		config: conf,
		stores: stores,
		debug:  true,
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

func (s *Server) initRateLimiter() error {
	conf := config.Get()

	if conf.API.RateLimitRequests <= 0 {
		slog.Warn("Rate limit requests is not set, rate limiting will be disabled")
	}

	lifetime := time.Minute
	rLimiter := NewRateLimiter(conf.API.RateLimitRequests, lifetime)
	s.rateLimit = rLimiter
	return nil
}

// Listen starts listening on the given address.
func (s *Server) Listen(addr string) error {

	if err := s.initAuthProvider(); err != nil {
		return err
	}

	if err := s.initRateLimiter(); err != nil {
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
	s.router.Use(middleware.RealIP)

	// auth middleware and routes
	s.router.Group(func(r chi.Router) {
		r.Use(authMiddleware(s.auth))
		r.Use(s.rateLimitMiddleware(s.config.API.RateLimitRequests, s.config.API.RateLimitStrategy))

		r.Post("/", makeAPIHandler(s.handlePost))
	})

	s.router.Get("/status", handleStatus)
	s.router.Get("/", makeAPIHandler(s.handleGet))

}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	status := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(status)
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) error {
	var req []*JRPCRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return writeError(w, http.StatusBadRequest, errReadingBody)
	}

	if len(body) == 0 {
		return writeError(w, http.StatusBadRequest, errMissingBody)
	}

	if err := json.Unmarshal(body, &req); err != nil {
		slog.Error("failed to unmarshal request", "err", err)
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

	var resp []Response
	// range over the requests and handle them
	for _, r := range req {
		response := s.handleJrpcRequest(r, authLevel)
		resp = append(resp, response)
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body></body></html>"))
	return nil
}

package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Llanos99/wf-engine/engine"
	"github.com/Llanos99/wf-engine/node"
	"github.com/Llanos99/wf-engine/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server is the REST API server
type Server struct {
	router       *chi.Mux
	store        store.Store
	executor     *engine.Executor
	registry     *node.Registry
	scriptRunner node.ScriptRunner
	httpServer   *http.Server
}

// Config holds server configuration
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DefaultConfig returns default server configuration
func DefaultConfig() Config {
	return Config{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
}

// NewServer creates a new API server
func NewServer(store store.Store, registry *node.Registry, scriptRunner node.ScriptRunner) *Server {
	s := &Server{
		router:       chi.NewRouter(),
		store:        store,
		executor:     engine.NewExecutor(registry, scriptRunner),
		registry:     registry,
		scriptRunner: scriptRunner,
	}
	s.setupRoutes()
	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	r := s.router

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", s.healthCheck)

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Definitions
		r.Route("/definitions", func(r chi.Router) {
			r.Get("/", s.listDefinitions)
			r.Post("/", s.createDefinition)
			r.Get("/{id}", s.getDefinition)
			r.Get("/{id}/versions/{version}", s.getDefinitionVersion)
			r.Delete("/{id}/versions/{version}", s.deleteDefinition)
		})

		// Instances
		r.Route("/instances", func(r chi.Router) {
			r.Get("/", s.listInstances)
			r.Post("/", s.createInstance)
			r.Get("/{id}", s.getInstance)
			r.Delete("/{id}", s.deleteInstance)
			r.Post("/{id}/resume", s.resumeInstance)
			r.Get("/{id}/logs", s.getInstanceLogs)
		})

		// Approvals
		r.Route("/approvals", func(r chi.Router) {
			r.Get("/", s.listPendingApprovals)
			r.Get("/{id}", s.getApproval)
			r.Post("/{id}/approve", s.approveRequest)
			r.Post("/{id}/reject", s.rejectRequest)
		})
	})
}

// Start starts the HTTP server
func (s *Server) Start(cfg Config) error {
	s.httpServer = &http.Server{
		Addr:         cfg.Addr,
		Handler:      s.router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	log.Printf("Starting API server on %s", cfg.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Router returns the chi router (for testing)
func (s *Server) Router() *chi.Mux {
	return s.router
}

// healthCheck handles GET /health
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// Helper to get URL parameters
func urlParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// Helper to parse int URL parameter
func urlParamInt(r *http.Request, key string) (int, error) {
	var val int
	_, err := fmt.Sscanf(chi.URLParam(r, key), "%d", &val)
	return val, err
}

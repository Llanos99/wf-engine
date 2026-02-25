package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Llanos99/wf-engine/api"
	"github.com/Llanos99/wf-engine/nodes"
	"github.com/Llanos99/wf-engine/script"
	"github.com/Llanos99/wf-engine/store/postgres"
)

func main() {
	// Flags
	addr := flag.String("addr", ":8080", "Server address")
	dbHost := flag.String("db-host", "localhost", "PostgreSQL host")
	dbPort := flag.Int("db-port", 5432, "PostgreSQL port")
	dbUser := flag.String("db-user", "postgres", "PostgreSQL user")
	dbPass := flag.String("db-pass", "postgres", "PostgreSQL password")
	dbName := flag.String("db-name", "wfengine", "PostgreSQL database")
	flag.Parse()

	// Database config
	dbCfg := postgres.Config{
		Host:     *dbHost,
		Port:     *dbPort,
		User:     *dbUser,
		Password: *dbPass,
		Database: *dbName,
		SSLMode:  "disable",
	}

	// Connect to database
	ctx := context.Background()
	store, err := postgres.NewStore(ctx, dbCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Close()
	log.Println("Connected to PostgreSQL")

	// Create node registry and script runner
	registry := nodes.NewRegistryWithBuiltin()
	scriptRunner := script.NewAdapter()

	// Create API server
	server := api.NewServer(store, registry, scriptRunner)

	// Server config
	serverCfg := api.Config{
		Addr:         *addr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start server in goroutine
	go func() {
		if err := server.Start(serverCfg); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

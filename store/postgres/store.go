package postgres

import (
	"context"
	"fmt"

	"github.com/Llanos99/wf-engine/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore implements store.Store using PostgreSQL
type PostgresStore struct {
	pool        *pgxpool.Pool
	definitions *DefinitionRepo
	instances   *InstanceRepo
	approvals   *ApprovalRepo
}

// Config holds the PostgreSQL connection configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// DefaultConfig returns a default configuration for local development
func DefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "wfengine",
		SSLMode:  "disable",
	}
}

// ConnectionString returns the PostgreSQL connection string
func (c Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

// NewStore creates a new PostgreSQL store with the given configuration
func NewStore(ctx context.Context, cfg Config) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{
		pool:        pool,
		definitions: NewDefinitionRepo(pool),
		instances:   NewInstanceRepo(pool),
		approvals:   NewApprovalRepo(pool),
	}, nil
}

// NewStoreFromPool creates a store from an existing connection pool
func NewStoreFromPool(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		pool:        pool,
		definitions: NewDefinitionRepo(pool),
		instances:   NewInstanceRepo(pool),
		approvals:   NewApprovalRepo(pool),
	}
}

// Definitions returns the definition repository
func (s *PostgresStore) Definitions() store.DefinitionRepository {
	return s.definitions
}

// Instances returns the instance repository
func (s *PostgresStore) Instances() store.InstanceRepository {
	return s.instances
}

// Approvals returns the approval repository
func (s *PostgresStore) Approvals() store.ApprovalRepository {
	return s.approvals
}

// Close closes the database connection pool
func (s *PostgresStore) Close() error {
	s.pool.Close()
	return nil
}

// Pool returns the underlying connection pool (for advanced use cases)
func (s *PostgresStore) Pool() *pgxpool.Pool {
	return s.pool
}

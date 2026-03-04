package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKey struct{}

// WithRepo adds repository to context
func WithRepo(ctx context.Context, repo *Repository) context.Context {
	return context.WithValue(ctx, ctxKey{}, repo)
}

// FromContext extracts repository from context
func FromContext(ctx context.Context) (*Repository, error) {
	repo, ok := ctx.Value(ctxKey{}).(*Repository)
	if !ok || repo == nil {
		return nil, errors.New("repository not found in context")
	}
	return repo, nil
}

// MustFromContext extracts repository from context or panics
func MustFromContext(ctx context.Context) *Repository {
	repo, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}
	return repo
}

type Repository struct {
	pool *pgxpool.Pool
}

// New creates a connection pool. Pass your DATABASE_URL here.
func New(ctx context.Context, databaseURL string) (*Repository, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	// Connections pool
	cfg.MaxConns = 10
	cfg.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	r := &Repository{pool: pool}

	if err := r.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return r, nil
}

func (r *Repository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *Repository) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return r.pool.Ping(ctx)
}

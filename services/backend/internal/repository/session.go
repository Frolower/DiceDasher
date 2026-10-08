package repository

import (
	"backend/internal/auth"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateSession(ctx context.Context, session auth.Session) error {
	const q = `INSERT INTO public.sessions (id, user_id, refresh_token_hash, created_at, expires_at, user_agent)
 VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, q, session.ID, session.UserID, session.RefreshTokenHash, session.CreatedAt, session.ExpiresAt, session.UserAgent)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

func (r *Repository) FindSessionByRefreshHash(ctx context.Context, hash string) (auth.Session, error) {
	const q = `SELECT id, user_id, refresh_token_hash, created_at, expires_at, revoked_at, COALESCE(user_agent, '') FROM public.sessions WHERE refresh_token_hash=$1`
	var s auth.Session
	err := r.pool.QueryRow(ctx, q, hash).Scan(&s.ID, &s.UserID, &s.RefreshTokenHash, &s.CreatedAt, &s.ExpiresAt, &s.RevokedAt, &s.UserAgent)
	if errors.Is(err, pgx.ErrNoRows) {
		return auth.Session{}, auth.ErrSessionNotFound
	}
	if err != nil {
		return auth.Session{}, fmt.Errorf("find session: %w", err)
	}
	return s, nil
}

func (r *Repository) RotateRefreshToken(ctx context.Context, sessionID uuid.UUID, oldHash, newHash string, now time.Time) error {
	tag, err := r.pool.Exec(ctx, `UPDATE public.sessions SET refresh_token_hash=$1 WHERE id=$2 AND refresh_token_hash=$3 AND revoked_at IS NULL AND expires_at>$4`, newHash, sessionID, oldHash, now)
	if err != nil {
		return fmt.Errorf("rotate session: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return auth.ErrSessionNotFound
	}
	return nil
}
func (r *Repository) RevokeSessionByRefreshHash(ctx context.Context, hash string, now time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE public.sessions SET revoked_at=$1 WHERE refresh_token_hash=$2 AND revoked_at IS NULL`, now, hash)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}

func (r *Repository) FindSessionByID(ctx context.Context, id uuid.UUID) (auth.Session, error) {
	const q = `SELECT id, user_id, refresh_token_hash, created_at, expires_at, revoked_at, COALESCE(user_agent, '') FROM public.sessions WHERE id=$1`
	var s auth.Session
	err := r.pool.QueryRow(ctx, q, id).Scan(&s.ID, &s.UserID, &s.RefreshTokenHash, &s.CreatedAt, &s.ExpiresAt, &s.RevokedAt, &s.UserAgent)
	if errors.Is(err, pgx.ErrNoRows) {
		return auth.Session{}, auth.ErrSessionNotFound
	}
	if err != nil {
		return auth.Session{}, fmt.Errorf("find session by id: %w", err)
	}

	return s, nil
}

package repository

import (
	"backend/internal/user"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetHashByID(ctx context.Context, id uuid.UUID) (string, error) {
	var hash string

	const q = `SELECT hash FROM public.users WHERE id = $1;`
	err := r.pool.QueryRow(ctx, q, id).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", user.ErrUserNotFound
	}
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (r *Repository) UpdatePasswordAndRevokeOtherSessions(
	ctx context.Context,
	userID, currentSessionID uuid.UUID,
	oldHash, newHash string,
	now time.Time,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin password update: %w", err)
	}

	defer func() {
		rollbackCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx), 5*time.Second,
		)
		defer cancel()
		_ = tx.Rollback(rollbackCtx)
	}()

	const updatePassword = `
		UPDATE public.users
		SET hash = $1
		WHERE id = $2 AND hash = $3
	`
	tag, err := tx.Exec(ctx, updatePassword, newHash, userID, oldHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return user.ErrUserNotFound
	}

	const revokeSessions = `
		UPDATE public.sessions
		SET revoked_at = $1
		WHERE user_id = $2
		  AND id <> $3
		  AND revoked_at IS NULL
	`
	if _, err := tx.Exec(
		ctx, revokeSessions, now, userID, currentSessionID,
	); err != nil {
		return fmt.Errorf("revoke other sessions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit password update: %w", err)
	}

	return nil
}

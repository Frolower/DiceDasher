package repository

import (
	"backend/internal/auth"
	"backend/internal/user"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) CreateUser(ctx context.Context, rec auth.CreateRecord) (uuid.UUID, error) {
	const q = `INSERT INTO public.users (username, email, hash)
               VALUES ($1, $2, $3) RETURNING id`
	var id uuid.UUID
	err := r.pool.QueryRow(ctx, q, rec.Username, rec.Email, rec.PasswordHash).Scan(&id)
	if err != nil {
		return uuid.Nil, mapCreateError(err)
	}
	return id, nil
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (auth.Credentials, error) {
	const q = `SELECT id, hash FROM public.users WHERE lower(username) = lower($1)`
	var credentials auth.Credentials
	err := r.pool.QueryRow(ctx, q, username).Scan(&credentials.UserID, &credentials.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return auth.Credentials{}, auth.ErrUserNotFound
	}
	if err != nil {
		return auth.Credentials{}, fmt.Errorf("find user: %w", err)
	}
	return credentials, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id uuid.UUID) (user.UserProfile, error) {
	const q = `SELECT id, username, email FROM public.users WHERE id = $1`
	var profile user.UserProfile
	err := r.pool.QueryRow(ctx, q, id).Scan(&profile.ID, &profile.Username, &profile.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.UserProfile{}, user.ErrUserNotFound
	}
	if err != nil {
		return user.UserProfile{}, fmt.Errorf("find user: %w", err)
	}
	return profile, nil
}

func (r *Repository) UpdateUser(ctx context.Context, u user.UserUpdateRecord) error {
	const q = `UPDATE public.users
	SET
    	username = COALESCE(NULLIF($1, ''), username),
    	email = COALESCE(NULLIF($2, ''), email),
		updated_at = NOW()
	WHERE id = $3;`
	tag, err := r.pool.Exec(ctx, q, u.Username, u.Email, u.ID)
	if err != nil {
		return mapUpdateError(err)
	}

	if tag.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM public.users WHERE id = $1;`
	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

func (r *Repository) UpdatePassword(ctx context.Context, id uuid.UUID, rec user.PasswordUpdateRecord) error {
	return nil
}

func mapCreateError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" && (pgErr.ConstraintName == "users_username_unique" || pgErr.ConstraintName == "users_email_unique") {
			return auth.ErrAlreadyExists
		}
		return fmt.Errorf("insert user: PostgreSQL error %s", pgErr.Code)
	}
	return fmt.Errorf("insert user: %w", err)
}

func mapUpdateError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" && (pgErr.ConstraintName == "users_username_unique" || pgErr.ConstraintName == "users_email_unique") {
			return user.ErrAlreadyExists
		}
		return fmt.Errorf("update user: PostgreSQL error %s", pgErr.Code)
	}
	return fmt.Errorf("update user: %w", err)
}

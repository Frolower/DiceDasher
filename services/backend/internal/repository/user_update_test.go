package repository

import (
	"backend/internal/user"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapUpdateError(t *testing.T) {
	for _, tc := range []struct {
		name, code, constraint string
		conflict               bool
	}{
		{"username conflict", "23505", "users_username_unique", true},
		{"email conflict", "23505", "users_email_unique", true},
		{"unrelated unique index", "23505", "users_pkey", false},
		{"other constraint violation", "23514", "users_username_unique", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := mapUpdateError(&pgconn.PgError{Code: tc.code, ConstraintName: tc.constraint, Message: "private message", Detail: "private email"})
			if errors.Is(err, user.ErrAlreadyExists) != tc.conflict {
				t.Fatalf("unexpected mapping: %v", err)
			}
			if strings.Contains(err.Error(), "private") {
				t.Fatal("leaked PostgreSQL details")
			}
		})
	}
	cause := errors.New("connection closed")
	if err := mapUpdateError(cause); !errors.Is(err, cause) {
		t.Fatalf("lost cause: %v", err)
	}
}

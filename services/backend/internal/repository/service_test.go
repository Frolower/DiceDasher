package repository

import (
	"backend/internal/auth"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"strings"
	"testing"
)

func TestMapCreateError(t *testing.T) {
	for _, name := range []string{"users_username_unique", "users_email_unique"} {
		if err := mapCreateError(&pgconn.PgError{Code: "23505", ConstraintName: name}); !errors.Is(err, auth.ErrAlreadyExists) {
			t.Fatal(err)
		}
	}
	// Other constraint violations must not be reported as an email or username already in use.
	err := mapCreateError(&pgconn.PgError{Code: "23505", ConstraintName: "users_pkey", Detail: "private", Message: "private"})
	if errors.Is(err, auth.ErrAlreadyExists) || strings.Contains(err.Error(), "private") {
		t.Fatal(err)
	}
	cause := errors.New("connection closed")
	if !errors.Is(mapCreateError(cause), cause) {
		t.Fatal("lost error cause")
	}
}

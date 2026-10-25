package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type updateStoreFunc struct {
	UserStore
	update func(context.Context, UserUpdateRecord) error
}

func (s updateStoreFunc) UpdateUser(ctx context.Context, rec UserUpdateRecord) error {
	return s.update(ctx, rec)
}

func TestUpdateUserNormalizesAndSaves(t *testing.T) {
	id := uuid.New()
	for _, tc := range []struct {
		name, username, email, wantUsername, wantEmail string
	}{
		{"both fields", " \tAlice\n", " Alice@Example.COM ", "Alice", "alice@example.com"},
		{"username only", "Alice", "", "Alice", ""},
		{"email only", "", "alice@example.com", "", "alice@example.com"},
		{"maximum ASCII username", strings.Repeat("a", 64), "", strings.Repeat("a", 64), ""},
		{"maximum email", "", strings.Repeat("a", 64) + "@" + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 61), "", strings.Repeat("a", 64) + "@" + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 61)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			calls := 0
			s := NewServiceWithDependencies(updateStoreFunc{update: func(got context.Context, rec UserUpdateRecord) error {
				calls++
				if got != ctx || rec != (UserUpdateRecord{ID: id, Username: tc.wantUsername, Email: tc.wantEmail}) {
					t.Fatalf("unexpected context or record: %+v", rec)
				}
				return nil
			}}, nil)
			if err := s.UpdateUser(ctx, UserUpdateRecord{ID: id, Username: tc.username, Email: tc.email}); err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Fatalf("store called %d times", calls)
			}
		})
	}
}

func TestUpdateUserInvalidInputDoesNotSave(t *testing.T) {
	id := uuid.New()
	for _, tc := range []struct {
		name    string
		input   UserUpdateRecord
		message string
	}{
		{"missing ID", UserUpdateRecord{Username: "Alice"}, "missing ID"},
		{"empty update", UserUpdateRecord{ID: id}, "missing arguments to update"},
		{"whitespace update", UserUpdateRecord{ID: id, Username: " \t", Email: "\n "}, "missing arguments to update"},
		{"short Unicode username", UserUpdateRecord{ID: id, Username: "Яяяя"}, "username must be between"},
		{"long Unicode username", UserUpdateRecord{ID: id, Username: strings.Repeat("я", 65)}, "username must be between"},
		{"internal Unicode whitespace", UserUpdateRecord{ID: id, Username: "Ali\u00a0ce"}, "username must not contain whitespace"},
		{"control character", UserUpdateRecord{ID: id, Username: "Ali\x00ce"}, "username must not contain control"},

		{"Cyrillic username", UserUpdateRecord{ID: id, Username: "Алиса"}, "username must use only English letters"},
		{"mixed alphabet username", UserUpdateRecord{ID: id, Username: "Alicе"}, "username must use only English letters"},
		{"accented username", UserUpdateRecord{ID: id, Username: "Alicé"}, "username must use only English letters"},
		{"emoji username", UserUpdateRecord{ID: id, Username: "Alice🎲"}, "username must use only English letters"},
		{"short email", UserUpdateRecord{ID: id, Email: "a@b"}, "email must be between"},
		{"long email", UserUpdateRecord{ID: id, Email: strings.Repeat("a", 243) + "@example.com"}, "email must be between"},
		{"malformed email", UserUpdateRecord{ID: id, Email: "not-an-email"}, "email is incorrect"},
		{"email without domain suffix", UserUpdateRecord{ID: id, Email: "frolowir@gmail"}, "email is incorrect"},
		{"valid username and invalid email", UserUpdateRecord{ID: id, Username: "Alice", Email: "invalid"}, "email is incorrect"},
		{"invalid username and valid email", UserUpdateRecord{ID: id, Username: "abc", Email: "alice@example.com"}, "username must be between"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := NewServiceWithDependencies(updateStoreFunc{update: func(context.Context, UserUpdateRecord) error {
				t.Fatal("invalid update reached store")
				return nil
			}}, nil)
			err := s.UpdateUser(context.Background(), tc.input)
			var invalid InvalidInputError
			if !errors.Is(err, ErrInvalidInput) || !errors.As(err, &invalid) || !strings.Contains(err.Error(), tc.message) {
				t.Fatalf("expected typed validation error containing %q, got %v", tc.message, err)
			}
		})
	}
}

func TestUpdateUserFailures(t *testing.T) {
	input := UserUpdateRecord{ID: uuid.New(), Username: "Alice"}
	if err := NewServiceWithDependencies(nil, nil).UpdateUser(context.Background(), input); !errors.Is(err, ErrStoreNotConfigured) {
		t.Fatalf("missing store: %v", err)
	}
	for _, cause := range []error{ErrAlreadyExists, ErrUserNotFound, context.Canceled, errors.New("database unavailable")} {
		t.Run(cause.Error(), func(t *testing.T) {
			failure := fmt.Errorf("store update: %w", cause)
			s := NewServiceWithDependencies(updateStoreFunc{update: func(context.Context, UserUpdateRecord) error { return failure }}, nil)
			if err := s.UpdateUser(context.Background(), input); !errors.Is(err, cause) {
				t.Fatalf("lost store error: %v", err)
			}
		})
	}
}

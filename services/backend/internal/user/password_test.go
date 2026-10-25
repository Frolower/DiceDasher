package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type passwordStoreStub struct {
	UserStore
	get    func(context.Context, uuid.UUID) (string, error)
	update func(context.Context, uuid.UUID, string, string) error
}

func (s passwordStoreStub) GetHashByID(ctx context.Context, id uuid.UUID) (string, error) {
	return s.get(ctx, id)
}

func (s passwordStoreStub) UpdatePassword(ctx context.Context, id uuid.UUID, oldHash, newHash string) error {
	return s.update(ctx, id, oldHash, newHash)
}

func TestUpdatePasswordWithBcrypt(t *testing.T) {
	hasher := BcryptHasher{Cost: bcrypt.MinCost}
	const oldPassword, newPassword = "OldPass1!", "NewPass2!"
	initialHash, err := hasher.Hash(oldPassword)
	if err != nil {
		t.Fatal(err)
	}
	id := uuid.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	savedHash := initialHash
	updates := 0
	store := passwordStoreStub{
		get: func(got context.Context, gotID uuid.UUID) (string, error) {
			if got != ctx || gotID != id {
				t.Fatal("lookup lost context or authenticated ID")
			}
			return savedHash, nil
		},
		update: func(got context.Context, gotID uuid.UUID, oldHash, newHash string) error {
			if got != ctx || gotID != id || oldHash != savedHash {
				t.Fatal("update must use the original stored hash, context and ID")
			}
			if newHash == newPassword || hasher.Verify(newHash, newPassword) != nil {
				t.Fatal("new password was not hashed correctly")
			}
			updates++
			savedHash = newHash
			return nil
		},
	}
	s := NewServiceWithDependencies(store, hasher)
	if err := s.UpdatePassword(ctx, PasswordUpdateRecord{Id: id, OldPassword: oldPassword, NewPassword: newPassword}); err != nil {
		t.Fatal(err)
	}
	if updates != 1 || savedHash == initialHash || hasher.Verify(savedHash, oldPassword) == nil {
		t.Fatal("password was not replaced")
	}
	// Reusing the old password must fail without another write.
	err = s.UpdatePassword(ctx, PasswordUpdateRecord{Id: id, OldPassword: oldPassword, NewPassword: "Another3!"})
	if !errors.Is(err, ErrOldPasswordDoesntMatch) || updates != 1 {
		t.Fatalf("reused password: err=%v updates=%d", err, updates)
	}
}

func TestUpdatePasswordRejectsInvalidNewPassword(t *testing.T) {
	for _, password := range []string{"", "Aa1!", strings.Repeat("a", 69) + "Aa1!", "Password1", "password1!", "PASSWORD1!", "Password!", "Pass word1!", "Password1!\t", "ПарольAa1!"} {
		t.Run(password, func(t *testing.T) {
			store := passwordStoreStub{get: func(context.Context, uuid.UUID) (string, error) {
				t.Fatal("invalid password reached storage")
				return "", nil
			}}
			err := NewServiceWithDependencies(store, nil).UpdatePassword(context.Background(), PasswordUpdateRecord{
				Id: uuid.New(), OldPassword: "OldPass1!", NewPassword: password,
			})
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	}
}

type passwordManagerStub struct {
	PasswordManager
	hash func(string) (string, error)
}

func (p passwordManagerStub) Hash(password string) (string, error) { return p.hash(password) }

func TestUpdatePasswordFailures(t *testing.T) {
	input := PasswordUpdateRecord{Id: uuid.New(), OldPassword: "OldPass1!", NewPassword: "NewPass2!"}
	if err := NewServiceWithDependencies(nil, nil).UpdatePassword(context.Background(), input); !errors.Is(err, ErrStoreNotConfigured) {
		t.Fatalf("missing store: %v", err)
	}
	hasher := BcryptHasher{Cost: bcrypt.MinCost}
	storedHash, err := hasher.Hash(input.OldPassword)
	if err != nil {
		t.Fatal(err)
	}
	for _, stage := range []string{"lookup missing user", "lookup database failure", "hash", "update"} {
		t.Run(stage, func(t *testing.T) {
			cause := errors.New("dependency failed")
			if stage == "lookup missing user" {
				cause = ErrUserNotFound
			}
			updates := 0
			store := passwordStoreStub{
				get: func(context.Context, uuid.UUID) (string, error) {
					if strings.HasPrefix(stage, "lookup") {
						return "", fmt.Errorf("lookup: %w", cause)
					}
					return storedHash, nil
				},
				update: func(context.Context, uuid.UUID, string, string) error {
					updates++
					return fmt.Errorf("update: %w", cause)
				},
			}
			manager := passwordManagerStub{PasswordManager: hasher, hash: func(password string) (string, error) {
				if stage == "hash" {
					return "", cause
				}
				return hasher.Hash(password)
			}}
			err := NewServiceWithDependencies(store, manager).UpdatePassword(context.Background(), input)
			if !errors.Is(err, cause) {
				t.Fatalf("lost cause: %v", err)
			}
			wantUpdates := 0
			if stage == "update" {
				wantUpdates = 1
			}
			if updates != wantUpdates {
				t.Fatalf("updates=%d want=%d", updates, wantUpdates)
			}
		})
	}
}

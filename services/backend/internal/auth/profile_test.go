package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"testing"
)

type profileStore struct {
	UserStore
	find func(context.Context, uuid.UUID) (UserProfile, error)
}

func (s profileStore) FindUserByID(ctx context.Context, id uuid.UUID) (UserProfile, error) {
	return s.find(ctx, id)
}
func TestGetProfile(t *testing.T) {
	failure := errors.New("database unavailable")
	for _, tc := range []struct {
		name string
		err  error
	}{
		{name: "success"},
		{name: "missing user", err: fmt.Errorf("find: %w", ErrUserNotFound)},
		{name: "database failure", err: failure},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want := UserProfile{ID: uuid.New(), Username: "Alice", Email: "alice@example.com"}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			calls := 0
			service := NewService(profileStore{find: func(received context.Context, id uuid.UUID) (UserProfile, error) {
				calls++
				if received != ctx || id != want.ID {
					t.Fatal("wrong lookup arguments")
				}
				if tc.err != nil {
					return UserProfile{}, tc.err
				}
				return want, nil
			}})
			got, err := service.GetProfile(ctx, want.ID)
			if !errors.Is(err, tc.err) || calls != 1 {
				t.Fatalf("err=%v calls=%d", err, calls)
			}
			if tc.err == nil && got != want {
				t.Fatalf("wrong profile: %v", got)
			}
			if tc.err != nil && got != (UserProfile{}) {
				t.Fatal("profile returned on failure")
			}
		})
	}
}
func TestGetProfileInvalidInput(t *testing.T) {
	service := NewService(profileStore{find: func(context.Context, uuid.UUID) (UserProfile, error) {
		t.Fatal("lookup with nil ID")
		return UserProfile{}, nil
	}})
	if got, err := service.GetProfile(context.Background(), uuid.Nil); !errors.Is(err, ErrUserNotFound) || got != (UserProfile{}) {
		t.Fatalf("got=%v err=%v", got, err)
	}
	if got, err := NewService(nil).GetProfile(context.Background(), uuid.New()); !errors.Is(err, ErrAuthNotConfigured) || got != (UserProfile{}) {
		t.Fatalf("got=%v err=%v", got, err)
	}
}

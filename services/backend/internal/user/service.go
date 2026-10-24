package user

import (
	"context"
	"time"
)

type UserStore interface {
	UpdateUser(ctx context.Context, rec UserUpdateRecord) error
	DeleteUser(ctx context.Context, id string) error
	UpdatePassword(ctx context.Context, id, hash string) error
}

type Service struct {
	users     UserStore
	passwords PasswordManager
	now       func() time.Time
}

// All dependencies are stored once when main wires up the application.
// A nil passwords dependency selects the default bcrypt implementation
func NewServiceWithDependencies(users UserStore, passwords PasswordManager) *Service {
	if passwords == nil {
		passwords = NewBcryptHasher()
	}
	return &Service{users: users, passwords: passwords, now: time.Now}
}

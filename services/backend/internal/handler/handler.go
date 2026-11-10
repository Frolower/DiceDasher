package handler

import (
	"backend/internal/auth"
	"backend/internal/user"
	"context"

	"github.com/google/uuid"
)

type Authenticator interface {
	Register(context.Context, auth.CreateInput) (auth.Created, error)
	Login(context.Context, auth.LoginInput) (auth.Tokens, error)
	Refresh(context.Context, string) (auth.Tokens, error)
	Logout(context.Context, string) error
	Authenticate(context.Context, string) (auth.AccessClaims, error)
}

type User interface {
	GetProfile(context.Context, uuid.UUID) (user.UserProfile, error)
	UpdateUser(ctx context.Context, record user.UserUpdateRecord) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, record user.PasswordUpdateRecord) error
}

type Rooms interface {
	Create(context.Context, string, uuid.UUID) (string, error)
}

type HealthChecker interface {
	Check(context.Context, string) error
}

type Options struct {
	CookieSecure   bool
	AllowedOrigins []string
}

type Handler struct {
	authenticator Authenticator
	user          User
	rooms         Rooms
	healthChecker HealthChecker
	options       Options
}

// By default, cookies are always secure, for local purposes it could be reconfigured in main

func New(authenticator Authenticator, users User) *Handler {
	return NewWithOptions(authenticator, users, nil, nil, Options{CookieSecure: true})
}

func NewWithOptions(authenticator Authenticator, user User, rooms Rooms, healthChecker HealthChecker, options Options) *Handler {
	return &Handler{
		authenticator: authenticator,
		user:          user,
		rooms:         rooms,
		healthChecker: healthChecker,
		options:       options,
	}
}

package handler

import (
	"backend/internal/auth"
	"context"
)

type Authenticator interface {
	Register(context.Context, auth.CreateInput) (auth.Created, error)
	Login(context.Context, auth.LoginInput) (auth.Tokens, error)
	Refresh(context.Context, string) (auth.Tokens, error)
	Logout(context.Context, string) error
}

type Options struct {
	CookieSecure  bool
	AllowedOrigin string
}
type Handler struct {
	authenticator Authenticator
	options       Options
}

// By default, cookies are always secure, for local purposes it could be reconfigured in main

func New(authenticator Authenticator) *Handler {
	return NewWithOptions(authenticator, Options{CookieSecure: true})
}
func NewWithOptions(authenticator Authenticator, options Options) *Handler {
	return &Handler{authenticator: authenticator, options: options}
}

package handler

import (
	"backend/internal/auth"
	"context"
)

type authContextKey struct{}

func withAccessClaims(ctx context.Context, claims auth.AccessClaims) context.Context {
	return context.WithValue(ctx, authContextKey{}, claims)
}

func AccessClaimsFromContext(ctx context.Context) (auth.AccessClaims, bool) {
	claims, ok := ctx.Value(authContextKey{}).(auth.AccessClaims)
	return claims, ok
}

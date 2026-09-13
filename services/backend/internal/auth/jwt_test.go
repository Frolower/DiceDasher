package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

func TestJWTRejectsInvalidTokens(t *testing.T) {
	key := []byte(strings.Repeat("k", 32))
	now := time.Now().UTC().Truncate(time.Second)
	m, _ := NewJWTManager(key, "issuer", "audience")
	m.now = func() time.Time { return now }
	for _, tc := range []struct {
		name   string
		modify func(*jwtClaims)
		method jwt.SigningMethod
		key    any
	}{
		{name: "valid"},
		{name: "expired", modify: func(c *jwtClaims) { c.ExpiresAt = jwt.NewNumericDate(now.Add(-time.Second)) }},
		{name: "missing exp", modify: func(c *jwtClaims) { c.ExpiresAt = nil }},
		{name: "wrong issuer", modify: func(c *jwtClaims) { c.Issuer = "other" }},
		{name: "missing audience", modify: func(c *jwtClaims) { c.Audience = nil }},
		{name: "wrong audience", modify: func(c *jwtClaims) { c.Audience = jwt.ClaimStrings{"other"} }},
		{name: "missing iat", modify: func(c *jwtClaims) { c.IssuedAt = nil }},
		{name: "future iat", modify: func(c *jwtClaims) { c.IssuedAt = jwt.NewNumericDate(now.Add(time.Hour)) }},
		{name: "future nbf", modify: func(c *jwtClaims) { c.NotBefore = jwt.NewNumericDate(now.Add(time.Hour)) }},
		{name: "invalid user", modify: func(c *jwtClaims) { c.Subject = "not-a-uuid" }},
		{name: "missing session", modify: func(c *jwtClaims) { c.SessionID = "" }},
		{name: "wrong key", key: []byte(strings.Repeat("z", 32))},
		{name: "wrong algorithm", method: jwt.SigningMethodHS384},
		{name: "unsigned", method: jwt.SigningMethodNone, key: jwt.UnsafeAllowNoneSignatureType},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := jwtClaims{SessionID: uuid.NewString(), RegisteredClaims: jwt.RegisteredClaims{Issuer: "issuer", Audience: jwt.ClaimStrings{"audience"}, Subject: uuid.NewString(), IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute))}}
			if tc.modify != nil {
				tc.modify(&c)
			}
			method := tc.method
			if method == nil {
				method = jwt.SigningMethodHS256
			}
			signingKey := tc.key
			if signingKey == nil {
				signingKey = key
			}
			raw, err := jwt.NewWithClaims(method, c).SignedString(signingKey)
			if err != nil {
				t.Fatal(err)
			}
			_, err = m.VerifyAccess(raw)
			if tc.name == "valid" {
				if err != nil {
					t.Fatal(err)
				}
			} else if !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("accepted invalid token: %v", err)
			}
		})
	}
	if _, err := m.VerifyAccess("broken"); !errors.Is(err, ErrInvalidToken) {
		t.Fatal(err)
	}
	if _, err := NewJWTManager([]byte("short"), "issuer", "audience"); err == nil {
		t.Fatal("accepted weak key")
	}
}

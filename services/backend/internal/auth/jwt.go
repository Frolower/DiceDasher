package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	key              []byte
	issuer, audience string
	now              func() time.Time
}
type jwtClaims struct {
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

var _ TokenManager = (*JWTManager)(nil)

func NewJWTManager(key []byte, issuer, audience string) (*JWTManager, error) {
	if len(key) < 32 || issuer == "" || audience == "" {
		return nil, errors.New("JWT requires a key of at least 32 bytes, issuer and audience")
	}
	return &JWTManager{key: append([]byte(nil), key...), issuer: issuer, audience: audience, now: time.Now}, nil
}

func (m *JWTManager) IssueAccess(userID, sessionID uuid.UUID, expiresAt time.Time) (string, error) {
	now := m.now().UTC()
	if userID == uuid.Nil || sessionID == uuid.Nil || !expiresAt.After(now) {
		return "", ErrInvalidToken
	}
	claims := jwtClaims{SessionID: sessionID.String(), RegisteredClaims: jwt.RegisteredClaims{
		Issuer: m.issuer, Subject: userID.String(), Audience: jwt.ClaimStrings{m.audience},
		IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(expiresAt), ID: uuid.NewString(),
	}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.key)
}

func (m *JWTManager) VerifyAccess(raw string) (AccessClaims, error) {
	if len(raw) > 8192 {
		return AccessClaims{}, ErrInvalidToken
	}
	var claims jwtClaims
	token, err := jwt.ParseWithClaims(raw, &claims, func(*jwt.Token) (any, error) { return m.key, nil },
		jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer(m.issuer), jwt.WithAudience(m.audience),
		jwt.WithExpirationRequired(), jwt.WithIssuedAt(), jwt.WithTimeFunc(m.now))
	if err != nil || !token.Valid || claims.IssuedAt == nil || claims.ExpiresAt == nil || !claims.ExpiresAt.After(claims.IssuedAt.Time) {
		return AccessClaims{}, ErrInvalidToken
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil || userID == uuid.Nil {
		return AccessClaims{}, ErrInvalidToken
	}
	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil || sessionID == uuid.Nil {
		return AccessClaims{}, ErrInvalidToken
	}
	return AccessClaims{UserID: userID, SessionID: sessionID, ExpiresAt: claims.ExpiresAt.Time}, nil
}

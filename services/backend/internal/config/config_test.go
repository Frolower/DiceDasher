package config

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestLoadJWTSettings(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("JWT_SECRET", base64.StdEncoding.EncodeToString([]byte(strings.Repeat("k", 32))))
	t.Setenv("COOKIE_SECURE", "")
	t.Setenv("ALLOWED_ORIGIN", "")
	c, err := Load()
	if err != nil || !c.CookieSecure || len(c.JWTKey) != 32 {
		t.Fatalf("config %v", err)
	}
	t.Setenv("COOKIE_SECURE", "false")
	c, err = Load()
	if err != nil || c.CookieSecure {
		t.Fatal("local override failed")
	}
	t.Setenv("JWT_SECRET", "short")
	if _, err = Load(); err == nil {
		t.Fatal("missing key validation")
	}
}

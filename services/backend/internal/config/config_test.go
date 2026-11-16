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
	t.Setenv("ALLOWED_ORIGINS", "")
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

func TestLoadHealthTargets(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("JWT_SECRET", base64.StdEncoding.EncodeToString([]byte(strings.Repeat("k", 32))))
	t.Setenv("COOKIE_SECURE", "")
	t.Setenv("ALLOWED_ORIGINS", "")
	t.Setenv("HEALTH_TARGETS", "")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.HealthTargets) != 3 || cfg.HealthTargets["room"] != "grpc://localhost:50051" {
		t.Fatalf("defaults: %v", cfg.HealthTargets)
	}
	t.Setenv("HEALTH_TARGETS", `{"another":"grpc://another:50051"}`)
	cfg, err = Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.HealthTargets) != 1 || cfg.HealthTargets["another"] != "grpc://another:50051" {
		t.Fatalf("override must replace defaults: %v", cfg.HealthTargets)
	}
	t.Setenv("HEALTH_TARGETS", `{}`)
	cfg, err = Load()
	if err != nil || len(cfg.HealthTargets) != 0 {
		t.Fatalf("disable targets: %v, %v", cfg, err)
	}
	for _, invalid := range []string{`null`, `[]`, `{"room":123}`, `not json`} {
		t.Setenv("HEALTH_TARGETS", invalid)
		if _, err := Load(); err == nil {
			t.Errorf("accepted %q", invalid)
		}
	}
}

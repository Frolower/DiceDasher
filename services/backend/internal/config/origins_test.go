package config

import (
	"encoding/base64"
	"reflect"
	"strings"
	"testing"
)

func TestLoadAllowedOrigins(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("ALLOWED_ORIGINS", "")
	t.Setenv("JWT_SECRET", base64.StdEncoding.EncodeToString([]byte(strings.Repeat("k", 32))))
	t.Setenv("COOKIE_SECURE", "false")

	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c.AllowedOrigins, []string{"http://localhost:8888", "http://localhost:5174"}) {
		t.Fatalf("unexpected defaults: %v", c.AllowedOrigins)
	}
	t.Setenv("ALLOWED_ORIGINS", " https://one.example, https://two.example ")
	c, err = Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c.AllowedOrigins, []string{"https://one.example", "https://two.example"}) {
		t.Fatalf("unexpected list: %v", c.AllowedOrigins)
	}
	t.Setenv("ALLOWED_ORIGINS", "https://one.example,https://two.example/path")
	if _, err := Load(); err == nil {
		t.Fatal("accepted invalid origin")
	}
}

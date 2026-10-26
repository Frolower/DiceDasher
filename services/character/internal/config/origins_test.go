package config

import (
	"reflect"
	"testing"
)

func TestLoadAllowedOrigins(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("ALLOWED_ORIGINS", "")

	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c.AllowedOrigins, []string{"http://localhost:8888", "http://localhost:8000"}) {
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

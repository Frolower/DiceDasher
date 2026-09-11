package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type storeFunc func(context.Context, CreateRecord) (uuid.UUID, error)

func (f storeFunc) CreateUser(ctx context.Context, rec CreateRecord) (uuid.UUID, error) {
	return f(ctx, rec)
}

type hashFunc func(string) (string, error)

func (f hashFunc) Hash(password string) (string, error) { return f(password) }

func TestRegisterNormalizesAndHashes(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	password := "StrongPass1"
	s := NewService(storeFunc(func(got context.Context, rec CreateRecord) (uuid.UUID, error) {
		if got != ctx || rec.Username != "Алиса" || rec.Email != "alice@example.com" {
			t.Fatalf("unexpected record: %s %s", rec.Username, rec.Email)
		}
		if rec.PasswordHash == password || bcrypt.CompareHashAndPassword([]byte(rec.PasswordHash), []byte(password)) != nil {
			t.Fatal("password was not hashed correctly")
		}
		return id, nil
	}))
	got, err := s.Register(ctx, CreateInput{Username: "Алиса", Email: " Alice@Example.COM ", Password: password})
	if err != nil || got.ID != id {
		t.Fatalf("result %v, error %v", got, err)
	}
}

func TestInvalidInputDoesNotHashOrSave(t *testing.T) {
	s := NewServiceWithHasher(storeFunc(func(context.Context, CreateRecord) (uuid.UUID, error) {
		t.Fatal("invalid input reached store")
		return uuid.Nil, nil
	}), hashFunc(func(string) (string, error) { t.Fatal("invalid input reached bcrypt"); return "", nil }))
	valid := CreateInput{Username: "Alice", Email: "alice@example.com", Password: "StrongPass1"}
	for _, tc := range []struct {
		name   string
		change func(*CreateInput)
	}{
		{"short Unicode username", func(v *CreateInput) { v.Username = "Яя" }},
		{"long username", func(v *CreateInput) { v.Username = strings.Repeat("я", 65) }},
		{"control in username", func(v *CreateInput) { v.Username = "Ali\x00ce" }},
		{"empty username", func(v *CreateInput) { v.Username = "   " }},
		{"bad email", func(v *CreateInput) { v.Email = "not-an-email" }},
		{"long email", func(v *CreateInput) { v.Email = strings.Repeat("a", 250) + "@example.com" }},
		{"short password", func(v *CreateInput) { v.Password = "Аа12345" }},
		{"missing uppercase", func(v *CreateInput) { v.Password = "password1" }},
		{"missing lowercase", func(v *CreateInput) { v.Password = "PASSWORD1" }},
		{"missing digit", func(v *CreateInput) { v.Password = "Password" }},
		{"bcrypt byte limit", func(v *CreateInput) { v.Password = "Aa1" + strings.Repeat("я", 35) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input := valid
			tc.change(&input)
			_, err := s.Register(context.Background(), input)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("error %v", err)
			}
		})
	}
}

func TestRegisterFailures(t *testing.T) {
	input := CreateInput{Username: "Alice", Email: "alice@example.com", Password: "StrongPass1"}
	failure := errors.New("hash failed")
	s := NewServiceWithHasher(storeFunc(func(context.Context, CreateRecord) (uuid.UUID, error) {
		t.Fatal("save after hash failure")
		return uuid.Nil, nil
	}), hashFunc(func(string) (string, error) { return "", failure }))
	if _, err := s.Register(context.Background(), input); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	for _, expected := range []error{ErrAlreadyExists, errors.New("storage failed")} {
		s = NewServiceWithHasher(storeFunc(func(context.Context, CreateRecord) (uuid.UUID, error) { return uuid.Nil, expected }), hashFunc(func(string) (string, error) { return "hash", nil }))
		if got, err := s.Register(context.Background(), input); !errors.Is(err, expected) || got.ID != uuid.Nil {
			t.Fatalf("result %v error %v", got, err)
		}
	}
	if _, err := NewService(nil).Register(context.Background(), input); !errors.Is(err, ErrStoreNotConfigured) {
		t.Fatal(err)
	}
}

func TestPasswordAtBcryptLimit(t *testing.T) {
	password := "Aa1" + strings.Repeat("x", 69)
	s := NewService(storeFunc(func(_ context.Context, rec CreateRecord) (uuid.UUID, error) {
		if err := bcrypt.CompareHashAndPassword([]byte(rec.PasswordHash), []byte(password)); err != nil {
			t.Fatal(err)
		}
		return uuid.New(), nil
	}))
	if _, err := s.Register(context.Background(), CreateInput{Username: "Alice", Email: "alice@example.com", Password: password}); err != nil {
		t.Fatal(err)
	}
}

// Проверяем сценарий целиком, чтобы нормализация не скрыла пробелы по краям.
func TestWhitespaceRejectedBeforeHashing(t *testing.T) {
	s := NewServiceWithHasher(storeFunc(func(context.Context, CreateRecord) (uuid.UUID, error) {
		t.Fatal("whitespace reached store")
		return uuid.Nil, nil
	}), hashFunc(func(string) (string, error) {
		t.Fatal("whitespace reached hasher")
		return "", nil
	}))
	for _, space := range []string{" ", "\t", "\n", "\r", "\u00a0", "\u2003"} {
		for _, position := range []string{"start", "middle", "end"} {
			for _, field := range []string{"username", "password"} {
				t.Run(fmt.Sprintf("%s/%s/%U", field, position, []rune(space)[0]), func(t *testing.T) {
					input := CreateInput{Username: "Alice", Email: "alice@example.com", Password: "StrongPass1"}
					value := &input.Username
					if field == "password" {
						value = &input.Password
					}
					switch position {
					case "start":
						*value = space + *value
					case "middle":
						*value = (*value)[:2] + space + (*value)[2:]
					case "end":
						*value += space
					}
					_, err := s.Register(context.Background(), input)
					if !errors.Is(err, ErrInvalidInput) || !strings.Contains(err.Error(), field+" must not contain whitespace") {
						t.Fatalf("unexpected error: %v", err)
					}
				})
			}
		}
	}
}

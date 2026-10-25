package handler

import (
	"backend/internal/auth"
	"backend/internal/repository"
	"backend/internal/user"
	"context"
	"diceDasher/pkg/httputil"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordUpdatePostgres(t *testing.T) {
	serviceURL, adminURL := os.Getenv("BACKEND_TEST_DATABASE_URL"), os.Getenv("BACKEND_TEST_ADMIN_URL")
	if serviceURL == "" || adminURL == "" {
		t.Skip("set BACKEND_TEST_DATABASE_URL and BACKEND_TEST_ADMIN_URL for a disposable database")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, serviceURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	admin, err := pgxpool.New(ctx, adminURL)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	store := repository.New(pool)
	name := "Password" + strings.ReplaceAll(uuid.NewString(), "-", "")
	const oldPassword, newPassword = "OldPass1!", "NewPass2!"
	created, err := auth.NewService(store).Register(ctx, auth.CreateInput{
		Username: name, Email: strings.ToLower(name) + "@example.com", Password: oldPassword,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.Exec(ctx, "DELETE FROM public.users WHERE id=$1", created.ID); err != nil {
			t.Error(err)
		}
	}()
	h := New(authenticateFunc{authenticate: func(context.Context, string) (auth.AccessClaims, error) {
		return auth.AccessClaims{UserID: created.ID, SessionID: uuid.New()}, nil
	}}, user.NewServiceWithDependencies(store, nil))
	router := httputil.NewRouter()
	h.RegisterRouters(router)
	request := func(old, next string) *httptest.ResponseRecorder {
		raw, err := json.Marshal(map[string]string{"old_password": old, "new_password": next})
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(http.MethodPut, "/password", strings.NewReader(string(raw)))
		req.Header.Set("Authorization", "Bearer test-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	readHash := func() string {
		var hash string
		if err := admin.QueryRow(ctx, "SELECT hash FROM public.users WHERE id=$1", created.ID).Scan(&hash); err != nil {
			t.Fatal(err)
		}
		return hash
	}
	initialHash := readHash()
	for _, tc := range []struct{ name, old, next string }{
		{"wrong old password", "WrongPass3!", newPassword},
		{"invalid new password", oldPassword, "weak"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if w := request(tc.old, tc.next); w.Code != http.StatusBadRequest {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if readHash() != initialHash {
				t.Fatal("rejected update changed the password")
			}
		})
	}
	if w := request(oldPassword, newPassword); w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	newHash := readHash()
	if newHash == initialHash || bcrypt.CompareHashAndPassword([]byte(newHash), []byte(newPassword)) != nil || bcrypt.CompareHashAndPassword([]byte(newHash), []byte(oldPassword)) == nil {
		t.Fatal("stored password was not replaced correctly")
	}
	if w := request(oldPassword, "Another3!"); w.Code != http.StatusBadRequest {
		t.Fatalf("old password accepted: %d", w.Code)
	}
	// A stale hash must not overwrite a password changed by another request.
	if err := store.UpdatePassword(ctx, created.ID, initialHash, initialHash); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("stale update: %v", err)
	}
	if readHash() != newHash {
		t.Fatal("stale update changed the password")
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := store.GetHashByID(canceled, created.ID); !errors.Is(err, context.Canceled) {
		t.Fatalf("lookup lost cancellation: %v", err)
	}
	if _, err := admin.Exec(ctx, "DELETE FROM public.users WHERE id=$1", created.ID); err != nil {
		t.Fatal(err)
	}
	if w := request(newPassword, "Another3!"); w.Code != http.StatusNotFound {
		t.Fatalf("deleted user: %d %s", w.Code, w.Body.String())
	}
}

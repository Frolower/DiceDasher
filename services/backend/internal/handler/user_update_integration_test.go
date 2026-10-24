package handler

import (
	"backend/internal/auth"
	"backend/internal/repository"
	"backend/internal/user"
	"context"
	"diceDasher/pkg/httputil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Exercise HTTP, validation, SQL, and the service role's permissions using a separate test database.
func TestUserUpdatePostgres(t *testing.T) {
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
	registration := auth.NewService(store)
	name := "Update" + strings.ReplaceAll(uuid.NewString(), "-", "")
	email := strings.ToLower(name) + "@example.com"
	created, err := registration.Register(ctx, auth.CreateInput{Username: name, Email: email, Password: "StrongPass1!"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.Exec(ctx, "DELETE FROM public.users WHERE id=$1", created.ID); err != nil {
			t.Error(err)
		}
	}()
	otherName, otherEmail := name+"Other", "other-"+email
	other, err := registration.Register(ctx, auth.CreateInput{Username: otherName, Email: otherEmail, Password: "StrongPass1!"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.Exec(ctx, "DELETE FROM public.users WHERE id=$1", other.ID); err != nil {
			t.Error(err)
		}
	}()
	h := New(authenticateFunc{authenticate: func(context.Context, string) (auth.AccessClaims, error) {
		return auth.AccessClaims{UserID: created.ID, SessionID: uuid.New()}, nil
	}}, user.NewServiceWithDependencies(store, nil))
	router := httputil.NewRouter()
	h.RegisterRouters(router)
	request := func(body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPatch, "/user", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer test-token")
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	readProfile := func(t *testing.T, id uuid.UUID) (string, string, *time.Time) {
		t.Helper()
		var username, address string
		var updatedAt *time.Time
		if err := admin.QueryRow(ctx, "SELECT username,email,updated_at FROM public.users WHERE id=$1", id).Scan(&username, &address, &updatedAt); err != nil {
			t.Fatal(err)
		}
		return username, address, updatedAt
	}
	wantName, wantEmail := name, email
	for _, tc := range []struct {
		name            string
		fields          map[string]string
		username, email string
		status          int
	}{
		{"username only", map[string]string{"username": name + "New"}, name + "New", email, 200},
		{"email only", map[string]string{"email": " NEW-" + strings.ToUpper(email) + " "}, name + "New", "new-" + email, 200},
		{"both fields", map[string]string{"username": name + "Both", "email": "both-" + email}, name + "Both", "both-" + email, 200},
		{"username conflict is atomic", map[string]string{"username": strings.ToLower(otherName), "email": "unused-" + email}, "", "", 409},
		{"email conflict is atomic", map[string]string{"username": name + "Unused", "email": strings.ToUpper(otherEmail)}, "", "", 409},
		{"empty update", map[string]string{}, "", "", 400},
		{"invalid email is atomic", map[string]string{"username": name + "Invalid", "email": "invalid"}, "", "", 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, before := readProfile(t, created.ID)
			raw, err := json.Marshal(tc.fields)
			if err != nil {
				t.Fatal(err)
			}
			w := request(string(raw))
			if w.Code != tc.status {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			if tc.status == 200 {
				wantName, wantEmail = tc.username, tc.email
			}
			gotName, gotEmail, after := readProfile(t, created.ID)
			if gotName != wantName || gotEmail != wantEmail {
				t.Fatalf("unexpected saved profile: %s %s", gotName, gotEmail)
			}
			if tc.status == 200 && after == nil {
				t.Fatal("updated_at was not set")
			}
			if tc.status != 200 && ((before == nil) != (after == nil) || (before != nil && after != nil && !before.Equal(*after))) {
				t.Fatal("failed update changed updated_at")
			}
			gotOtherName, gotOtherEmail, _ := readProfile(t, other.ID)
			if gotOtherName != otherName || gotOtherEmail != otherEmail {
				t.Fatal("changed another user's profile")
			}
		})
	}
	// A valid authenticated ID with no user row must not produce a successful update.
	if _, err := admin.Exec(ctx, "DELETE FROM public.users WHERE id=$1", created.ID); err != nil {
		t.Fatal(err)
	}
	if w := request(`{"username":"Alice"}`); w.Code != http.StatusNotFound {
		t.Fatalf("deleted user: %d %s", w.Code, w.Body.String())
	}
}

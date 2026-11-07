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
	"time"

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
	currentSessionID := uuid.New()
	otherSessionID, revokedSessionID, foreignSessionID := uuid.New(), uuid.New(), uuid.New()
	foreignUserID := uuid.New()
	if _, err := admin.Exec(ctx, `INSERT INTO public.users (id, username, email, hash) VALUES ($1, $2, $3, $4)`, foreignUserID, name+"Other", strings.ToLower(name)+"other@example.com", "unused"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.Exec(ctx, "DELETE FROM public.users WHERE id=$1", foreignUserID); err != nil {
			t.Error(err)
		}
	}()
	previousRevocation := time.Now().UTC().Add(-time.Hour).Truncate(time.Microsecond)
	for _, session := range []struct {
		id, userID uuid.UUID
		revokedAt  *time.Time
	}{
		{currentSessionID, created.ID, nil},
		{otherSessionID, created.ID, nil},
		{revokedSessionID, created.ID, &previousRevocation},
		{foreignSessionID, foreignUserID, nil},
	} {
		if _, err := admin.Exec(ctx, `INSERT INTO public.sessions (id, user_id, refresh_token_hash, revoked_at) VALUES ($1, $2, $3, $4)`, session.id, session.userID, session.id.String(), session.revokedAt); err != nil {
			t.Fatal(err)
		}
	}
	readRevocation := func(id uuid.UUID) *time.Time {
		var revokedAt *time.Time
		if err := admin.QueryRow(ctx, "SELECT revoked_at FROM public.sessions WHERE id=$1", id).Scan(&revokedAt); err != nil {
			t.Fatal(err)
		}
		return revokedAt
	}
	h := New(authenticateFunc{authenticate: func(context.Context, string) (auth.AccessClaims, error) {
		return auth.AccessClaims{UserID: created.ID, SessionID: currentSessionID}, nil
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
	readUpdatedAt := func() *time.Time {
		var updatedAt *time.Time
		if err := admin.QueryRow(ctx, "SELECT updated_at FROM public.users WHERE id=$1", created.ID).Scan(&updatedAt); err != nil {
			t.Fatal(err)
		}
		return updatedAt
	}
	assertUpdatedAt := func(want *time.Time) {
		t.Helper()
		got := readUpdatedAt()
		if (got == nil) != (want == nil) || (got != nil && want != nil && !got.Equal(*want)) {
			t.Fatalf("updated_at changed unexpectedly: got=%v want=%v", got, want)
		}
	}
	initialHash := readHash()
	initialUpdatedAt := readUpdatedAt()
	for _, tc := range []struct{ name, old, next string }{
		{"wrong old password", "WrongPass3!", newPassword},
		{"invalid new password", oldPassword, "weak"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if w := request(tc.old, tc.next); w.Code != http.StatusBadRequest {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			assertUpdatedAt(initialUpdatedAt)
			if readRevocation(otherSessionID) != nil {
				t.Fatal("rejected update revoked a session")
			}
			if readHash() != initialHash {
				t.Fatal("rejected update changed the password")
			}
		})
	}
	// Use the database clock to avoid clock skew between the app and PostgreSQL.
	var beforeUpdate time.Time
	if err := admin.QueryRow(ctx, "SELECT clock_timestamp()").Scan(&beforeUpdate); err != nil {
		t.Fatal(err)
	}
	if w := request(oldPassword, newPassword); w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if readRevocation(currentSessionID) != nil || readRevocation(foreignSessionID) != nil {
		t.Fatal("current session or another user's session was revoked")
	}
	if readRevocation(otherSessionID) == nil {
		t.Fatal("other session was not revoked")
	}
	if got := readRevocation(revokedSessionID); got == nil || !got.Equal(previousRevocation) {
		t.Fatal("previous revocation timestamp was changed")
	}
	updatedAt := readUpdatedAt()
	var afterUpdate time.Time
	if err := admin.QueryRow(ctx, "SELECT clock_timestamp()").Scan(&afterUpdate); err != nil {
		t.Fatal(err)
	}
	if updatedAt == nil || updatedAt.Before(beforeUpdate) || updatedAt.After(afterUpdate) {
		t.Fatalf("updated_at must reflect the successful password update: got=%v interval=[%v, %v]", updatedAt, beforeUpdate, afterUpdate)
	}
	newHash := readHash()
	if newHash == initialHash || bcrypt.CompareHashAndPassword([]byte(newHash), []byte(newPassword)) != nil || bcrypt.CompareHashAndPassword([]byte(newHash), []byte(oldPassword)) == nil {
		t.Fatal("stored password was not replaced correctly")
	}
	if w := request(oldPassword, "Another3!"); w.Code != http.StatusBadRequest {
		t.Fatalf("old password accepted: %d", w.Code)
	}
	assertUpdatedAt(updatedAt)
	// A stale hash must not overwrite a password changed by another request.
	if err := store.UpdatePasswordAndRevokeOtherSessions(ctx, created.ID, currentSessionID, initialHash, initialHash, time.Now().UTC()); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("stale update: %v", err)
	}
	assertUpdatedAt(updatedAt)
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

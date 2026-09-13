package handler

import (
	"backend/internal/auth"
	"backend/internal/repository"
	"context"
	"diceDasher/pkg/httputil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Весь путь HTTP → auth → SQL проверяется с реальными JWT и bcrypt.
// Рабочую БД тест не использует: URL отдельной тестовой БД задаётся явно.
func TestLoginPostgres(t *testing.T) {
	serviceURL, adminURL := os.Getenv("BACKEND_TEST_DATABASE_URL"), os.Getenv("BACKEND_TEST_ADMIN_URL")
	if serviceURL == "" || adminURL == "" {
		t.Skip("set disposable database URLs")
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
	jwtManager, err := auth.NewJWTManager([]byte(strings.Repeat("k", 32)), "test-issuer", "test-api")
	if err != nil {
		t.Fatal(err)
	}
	service := auth.NewServiceWithDependencies(store, store, auth.NewBcryptHasher(), jwtManager)
	name := "Login" + strings.ReplaceAll(uuid.NewString(), "-", "")
	created, err := service.Register(ctx, auth.CreateInput{Username: name, Email: name + "@example.com", Password: "StrongPass1"})
	if err != nil {
		t.Fatal(err)
	}
	// ON DELETE CASCADE удалит только сессии созданного тестом пользователя.
	defer func() {
		if _, err := admin.Exec(ctx, "DELETE FROM public.users WHERE id=$1", created.ID); err != nil {
			t.Error(err)
		}
	}()
	router := httputil.NewRouter()
	NewWithOptions(service, Options{CookieSecure: false, AllowedOrigin: "http://localhost:8082"}).RegisterRouters(router)
	request := func(username, password string) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(loginRequest{Username: username, Password: password})
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(raw)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "integration browser")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		return w
	}
	var firstRefresh string
	for i := 0; i < 2; i++ {
		w := request(strings.ToLower(name), "StrongPass1")
		if w.Code != 200 {
			t.Fatalf("login %d: %s", w.Code, w.Body.String())
		}
		var response tokenResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		claims, err := jwtManager.VerifyAccess(response.AccessToken)
		if err != nil || claims.UserID != created.ID || !claims.ExpiresAt.Equal(response.ExpiresAt) {
			t.Fatalf("invalid access token: %v", err)
		}
		cookies := w.Result().Cookies()
		if len(cookies) != 1 {
			t.Fatal("missing refresh cookie")
		}
		refresh := cookies[0].Value
		if refresh == firstRefresh {
			t.Fatal("reused refresh token")
		}
		firstRefresh = refresh
		if strings.Contains(w.Body.String(), refresh) {
			t.Fatal("refresh leaked to JSON")
		}
		session, err := store.FindSessionByRefreshHash(ctx, jwtManager.HashRefresh(refresh))
		if err != nil {
			t.Fatal(err)
		}
		if session.ID != claims.SessionID || session.UserID != created.ID || session.RefreshTokenHash == refresh || session.UserAgent != "integration browser" || session.RevokedAt != nil || !session.ExpiresAt.Equal(cookies[0].Expires) {
			t.Fatalf("incorrect stored session: %v", session.ID)
		}
	}
	bad := request(name, "WrongPass1")
	missing := request(name+"Missing", "StrongPass1")
	if bad.Code != 401 || missing.Code != 401 || bad.Body.String() != missing.Body.String() || len(bad.Result().Cookies()) != 0 || len(missing.Result().Cookies()) != 0 {
		t.Fatal("credential failures reveal user or return cookie")
	}
	var count int
	if err := admin.QueryRow(ctx, "SELECT count(*) FROM public.sessions WHERE user_id=$1", created.ID).Scan(&count); err != nil || count != 2 {
		t.Fatalf("unexpected sessions: %d %v", count, err)
	}
	// Если PostgreSQL недоступен, клиент не получает токены или cookie.
	pool.Close()
	w := request(name, "StrongPass1")
	if w.Code != 500 || len(w.Result().Cookies()) != 0 || strings.Contains(w.Body.String(), "pool") {
		t.Fatal("incorrect storage failure response")
	}
}

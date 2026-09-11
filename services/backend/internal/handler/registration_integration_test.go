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
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func TestRegistrationPostgres(t *testing.T) {
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
	name := "Test" + strings.ReplaceAll(uuid.NewString(), "-", "")
	email := strings.ToLower(name) + "@example.com"
	raceName := name + "Race"
	defer func() {
		if _, err := admin.Exec(ctx, "DELETE FROM public.users WHERE username = $1 OR username = $2", name, raceName); err != nil {
			t.Error(err)
		}
	}()
	router := httputil.NewRouter()
	New(auth.NewService(repository.New(pool))).RegisterRouters(router)
	request := func(username, address, password string) *httptest.ResponseRecorder {
		body, err := json.Marshal(createUserRequest{Username: username, Email: address, Password: password})
		if err != nil {
			panic(err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(body))))
		return w
	}
	w := request(name, " "+strings.ToUpper(email)+" ", "StrongPass1")
	if w.Code != 201 || w.Body.String() != "{\"status\":\"created\"}\n" || len(w.Result().Cookies()) != 0 {
		t.Fatalf("registration: %d %s", w.Code, w.Body.String())
	}
	var savedEmail, hash string
	if err := admin.QueryRow(ctx, "SELECT email, hash FROM public.users WHERE username=$1", name).Scan(&savedEmail, &hash); err != nil {
		t.Fatal(err)
	}
	if savedEmail != email || bcrypt.CompareHashAndPassword([]byte(hash), []byte("StrongPass1")) != nil {
		t.Fatal("incorrect normalized email or password hash")
	}
	for _, input := range [][2]string{{strings.ToLower(name), "other-" + email}, {name + "Other", strings.ToUpper(email)}} {
		if w := request(input[0], input[1], "StrongPass1"); w.Code != 409 {
			t.Fatalf("duplicate: %d %s", w.Code, w.Body.String())
		}
	}
	if w := request(name+"Invalid", "invalid-"+email, "short"); w.Code != 400 {
		t.Fatalf("validation: %d", w.Code)
	}
	statuses := make(chan int, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; statuses <- request(raceName, "race-"+email, "StrongPass1").Code }()
	}
	close(start)
	wg.Wait()
	close(statuses)
	counts := map[int]int{}
	for status := range statuses {
		counts[status]++
	}
	if counts[201] != 1 || counts[409] != 1 {
		t.Fatalf("concurrent registrations: %v", counts)
	}
	var count int
	if err := admin.QueryRow(ctx, "SELECT count(*) FROM public.users WHERE email = $1 OR email = $2", email, "race-"+email).Scan(&count); err != nil || count != 2 {
		t.Fatalf("stored count %d: %v", count, err)
	}
}

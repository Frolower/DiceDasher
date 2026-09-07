package main

import (
	"context"
	"diceDasher/pkg/dbutil"
	"diceDasher/pkg/httputil"
	"diceDasher/services/character/internal/config"
	"diceDasher/services/character/internal/handler"
	"diceDasher/services/character/internal/repository"
	"diceDasher/services/character/internal/service"
	"diceDasher/services/character/internal/system"
	"diceDasher/services/character/internal/system/tes"
	"log"
	"net"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := dbutil.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	characters := service.New(repository.New(repo.Pool()), map[string]system.Character{"tes": tes.Character{}})

	r := httputil.NewRouter()
	handler.New(characters).RegisterRouters(r)

	// Wrap with middlewares
	wrapped := httputil.RequestLoggerWithMode(r, cfg.LogMode)
	wrapped = httputil.CORS("http://localhost:8082")(wrapped)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: wrapped,
	}

	ln, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("character service READY on %s", cfg.HTTPAddr)
	log.Fatal(srv.Serve(ln))
}

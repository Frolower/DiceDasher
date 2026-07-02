package main

import (
	"context"
	"diceDasher/pkg/dbutil"
	"log"
	"net"
	"net/http"

	"diceDasher/pkg/httputil"
	"diceDasher/services/resolve/internal/config"
	"diceDasher/services/resolve/internal/handler"
	"diceDasher/services/resolve/internal/system"
	"diceDasher/services/resolve/internal/system/generic"
	"diceDasher/services/resolve/internal/system/tes"
	"diceDasher/services/resolve/internal/system/vtmv5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repository
	repo, err := dbutil.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	// register resolvers
	system.Register("generic", generic.Resolver{})
	system.Register("tes", tes.Resolver{})
	system.Register("vtmv5", vtmv5.Resolver{})

	r := httputil.NewRouter()
	handler.RegisterRoutes(r)

	// Wrap with middlewares
	wrapped := handler.WithRepository(repo)(r)
	wrapped = httputil.RequestLoggerWithMode(wrapped, "default")
	wrapped = httputil.CORS("http://localhost:8082")(wrapped)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: wrapped,
	}

	ln, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("resolver service READY on %s", cfg.HTTPAddr)
	log.Fatal(srv.Serve(ln))
}

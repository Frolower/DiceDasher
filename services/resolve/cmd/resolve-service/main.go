package main

import (
	"context"
	"diceDasher/pkg/dbutil"
	"diceDasher/pkg/dice"
	"diceDasher/pkg/logger"
	"diceDasher/services/resolve/internal/repository"
	"diceDasher/services/resolve/internal/service"
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
		logger.LogFatalf(context.Background(), "%v", err)
	}

	// Initialize repository
	repo, err := dbutil.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.LogFatalf(context.Background(), "%v", err)
	}
	defer repo.Close()

	history := repository.New(repo.Pool())
	generator := dice.Generator{}
	resolveService := service.New(map[string]system.Resolver{
		"generic": generic.Resolver{Dice: generator},
		"tes":     tes.New(history, generator),
		"vtmv5":   vtmv5.New(history, generator),
	}, history)
	r := httputil.NewRouter()
	handler.New(resolveService).RegisterRoutes(r)
	wrapped := httputil.RequestLoggerWithMode(r, cfg.LogMode)
	wrapped = httputil.CORS("http://localhost:8081")(wrapped)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: wrapped,
	}

	ln, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		logger.LogFatalf(context.Background(), "%v", err)
	}

	logger.Logf(context.Background(), "resolver service READY on %s", cfg.HTTPAddr)
	logger.LogFatalf(context.Background(), "%v", srv.Serve(ln))
}

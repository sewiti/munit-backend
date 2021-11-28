package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/sewiti/munit-backend/internal/auth"
	"github.com/sewiti/munit-backend/internal/config"
	"github.com/sewiti/munit-backend/internal/model"
	"github.com/sewiti/munit-backend/internal/web"
	"github.com/vrischmann/envconfig"
)

func main() {
	var cfg struct{ Munit config.Munit }
	if err := envconfig.Init(&cfg); err != nil {
		log.WithError(err).Fatal("unable to read environment config")
		return
	}

	if cfg.Munit.Debug {
		log.SetLevelFromString("debug")
	}

	if err := auth.LoadSecret(cfg.Munit.SecretFile); err != nil {
		log.WithError(err).Fatal("unable to setup secret")
		return
	}

	if err := model.OpenDB(cfg.Munit.DSN); err != nil {
		log.WithError(err).Fatal("unable to open database")
		return
	}
	defer model.CloseDB()

	// Create server
	srv := &http.Server{
		Addr:         cfg.Munit.Addr,
		Handler:      web.NewRouter(&cfg.Munit),
		WriteTimeout: cfg.Munit.Timeout,
		ReadTimeout:  cfg.Munit.Timeout,
	}

	// Start server
	go func() {
		log.Infof("server started on %s", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("server listen and serve")
		}
	}()

	// Wait for shutdown signal
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()
	log.Info("shutting down")

	// Initiate graceful server shutdown with a timeout
	ctx, cancel = context.WithTimeout(context.Background(), cfg.Munit.Timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Error("server shutdown")
	}
}

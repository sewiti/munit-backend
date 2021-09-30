package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/vrischmann/envconfig"
)

var conf struct {
	Munit struct {
		Addr          string        `envconfig:"default=:7878"`
		AllowedOrigin string        `envconfig:"default=munit.digital"`
		Timeout       time.Duration `envconfig:"default=30s"`
	}
}

func main() {
	if err := envconfig.Init(&conf); err != nil {
		log.WithError(err).Fatal("unable to read environment config")
		return
	}

	r := mux.NewRouter()

	// Setup CORS
	origins := handlers.AllowedOrigins([]string{conf.Munit.AllowedOrigin})
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete})

	// Create server
	srv := &http.Server{
		Addr:         conf.Munit.Addr,
		Handler:      handlers.CORS(origins, headers, methods)(r),
		WriteTimeout: conf.Munit.Timeout,
		ReadTimeout:  conf.Munit.Timeout,
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

	// Initiate gracefull server shutdown with a timeout
	ctx, cancel = context.WithTimeout(context.Background(), conf.Munit.Timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Error("server shutdown")
	}
}

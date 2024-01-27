package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jwtly10/simple-site-tracker/api/middleware"
	. "github.com/jwtly10/simple-site-tracker/api/router"
	"github.com/jwtly10/simple-site-tracker/api/service"
	"github.com/jwtly10/simple-site-tracker/api/track"
	"github.com/jwtly10/simple-site-tracker/config"
	"github.com/jwtly10/simple-site-tracker/utils/logger"
)

func main() {
	l := logger.Get()

	cfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal().Err(err).Msg("Error loading config")
	}

	db, err := config.OpenDB(cfg)
	if err != nil {
		l.Fatal().Err(err).Msg("Error opening database connection")
	}
	defer db.Close()

	// Load repository and handlers
	repo := track.NewRepository(db)
	th := track.NewHandlers(repo)

	svc := service.NewService(repo)
	mw := middleware.NewMiddleware(svc)

	router := NewRouter(th, mw)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		l.Info().Msg("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal().Err(err).Msg("Error starting server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	l.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		l.Fatal().Err(err).Msg("Error shutting down server")
	}

	l.Info().Msg("Server gracefully stopped")
}

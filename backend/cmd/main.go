package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"song-match-backend/api/route"
	"song-match-backend/bootstrap"
	"syscall"
	"time"
)

func main() {
	app := bootstrap.App()
	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second
	router, trackController := route.Setup(env, timeout, db)

	srv := &http.Server{
		Addr:    env.ServerAddress,
		Handler: router,
	}

	// Listen for SIGINT/SIGTERM and perform graceful shutdown.
	// This ensures background goroutines finish before the process exits,
	// preventing tracks from being permanently stuck in "processing".
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("starting HTTP server", "address", env.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	// Give in-flight HTTP requests 15 seconds to complete.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error", "error", err)
	}

	// Wait for all background processing goroutines to finish.
	slog.Info("waiting for background processing to complete")
	trackController.TrackUsecase.Shutdown()
	slog.Info("shutdown complete")
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"bff-go-mvp/internal/config"
	"bff-go-mvp/internal/logger"
	"bff-go-mvp/internal/router"
)

func main() {
	// Initialize logger
	zapLogger, err := logger.NewLogger(os.Getenv("ENV"))
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer zapLogger.Sync()

	// Load configuration
	cfg := config.Load()

	// Setup router with all endpoints
	r := router.New(cfg, zapLogger)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.API.Port)
	zapLogger.Info("Starting API server", zap.String("address", serverAddr))

	// Setup graceful shutdown
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zapLogger.Info("Server exited")
}

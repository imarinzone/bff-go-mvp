package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"

	"bff-go-mvp/internal/api"
	"bff-go-mvp/internal/config"
	"bff-go-mvp/internal/logger"
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

	// Create Temporal client
	temporalClient, err := client.NewClient(client.Options{
		HostPort:  cfg.Temporal.Host,
		Namespace: cfg.Temporal.Namespace,
	})
	if err != nil {
		zapLogger.Fatal("Unable to create Temporal client", zap.Error(err))
	}
	defer temporalClient.Close()

	// Create discovery handler (temporalClient implements TemporalWorkflowExecutor interface)
	discoveryHandler := api.NewDiscoveryHandler(
		temporalClient,
		cfg.GRPC.ServiceAddress,
		cfg.Temporal.TaskQueue,
		zapLogger,
	)

	// Setup router
	r := mux.NewRouter()
	
	// Middleware for logging
	r.Use(loggingMiddleware(zapLogger))
	r.Use(recoveryMiddleware(zapLogger))

	// Register routes
	r.HandleFunc("/discovery", discoveryHandler.HandleDiscovery).Methods("POST")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

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

// loggingMiddleware logs HTTP requests
func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", duration),
			)
		})
	}
}

// recoveryMiddleware recovers from panics
func recoveryMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						zap.Any("error", err),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
					)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

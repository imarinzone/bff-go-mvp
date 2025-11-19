package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"

	"bff-go-mvp/internal/config"
	"bff-go-mvp/internal/logger"
	"bff-go-mvp/internal/temporal"
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

	// Create worker
	w := worker.New(temporalClient, cfg.Temporal.TaskQueue, worker.Options{})

	// Register workflow
	w.RegisterWorkflow(temporal.DiscoveryWorkflow)

	// Register activity
	w.RegisterActivity(temporal.CallDiscoveryActivity)

	// Start worker
	zapLogger.Info("Starting Temporal worker",
		zap.String("task_queue", cfg.Temporal.TaskQueue),
		zap.String("namespace", cfg.Temporal.Namespace),
	)

	// Start worker in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			errChan <- err
		}
	}()

	// Wait for interrupt signal or error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		zapLogger.Info("Shutting down worker...")
		w.Stop()
		zapLogger.Info("Worker stopped")
	case err := <-errChan:
		zapLogger.Fatal("Worker failed", zap.Error(err))
	}
}


package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"conduit/internal/config"
	"conduit/internal/loadbalancer"
	"conduit/internal/server"
	"conduit/internal/servicemesh"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize service mesh
	sm := servicemesh.NewServiceMesh(logger)
	go sm.HealthCheck()

	// Initialize load balancer
	lb, err := loadbalancer.NewLoadBalancer(cfg.BackendURLs, logger)
	if err != nil {
		logger.Fatal("Failed to create load balancer", zap.Error(err))
	}
	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			lb.HealthCheck()
		}
	}()

	// Initialize server
	s, err := server.NewServer(cfg, logger, sm)
	if err != nil {
		logger.Fatal("Failed to create server", zap.Error(err))
	}

	// Register this server with the service mesh
	sm.RegisterService("quant-trading", "http://"+cfg.ServerAddress)

	// Set up HTTP server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: lb,
	}

	// Run our server in a goroutine so that it doesn't block
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("ListenAndServe error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	sm.DeregisterService("quant-trading", "http://"+cfg.ServerAddress)
	logger.Info("Server exiting")
}

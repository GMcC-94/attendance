package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gmcc94/attendance-go/internal/config"
	"github.com/gmcc94/attendance-go/internal/db"
	"github.com/gmcc94/attendance-go/internal/router"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Init DB
	database, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialise database: %w", err)
	}
	defer database.Close()

	// Setup router
	r := router.New(database, cfg)

	// Create server
	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server starting on %s", cfg.ServerAddress)
	log.Printf("🌍 Environment: %s", cfg.Environment)

	// Start server with graceful shutdown
	return startServerWithGracefulShutdown(server)
}

func startServerWithGracefulShutdown(server *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-quit
	log.Println("🛑 Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("✅ Server exited")
	return nil
}

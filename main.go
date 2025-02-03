package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmarquinah/publist_backend/internal/config"
	"github.com/dmarquinah/publist_backend/internal/handler"
	"github.com/dmarquinah/publist_backend/internal/repository"
	"github.com/dmarquinah/publist_backend/internal/service"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize dependencies
	repo := repository.NewRepository()
	svc := service.NewService(repo)
	handlers := handler.NewHandler(svc)

	// Setup router
	mux := http.NewServeMux()

	// Register routes
	handlers.RegisterRoutes(mux)

	// Configure server
	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}

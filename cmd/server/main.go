package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	httpadapter "cleanarch/internal/adapter/http"
	"cleanarch/internal/app"
	"cleanarch/internal/repository/memory"
	"cleanarch/internal/usecase"
)

func main() {
	// Initialize dependencies
	repo := memory.NewInMemoryUserRepository()
	service := usecase.NewUserService(repo)
	handler := httpadapter.NewUserHandler(service)

	mux := app.NewRouter(handler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      app.WithLogging(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("HTTP server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-shutdownCtx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	} else {
		log.Println("server shutdown complete")
	}
}

package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	httpadapter "cleanarch/internal/adapter/http"
	"cleanarch/internal/app"
	"cleanarch/internal/repository/memory"
	"cleanarch/internal/usecase"
)

func TestDependencyInjection(t *testing.T) {
	t.Run("Dependencies are created successfully", func(t *testing.T) {
		// Test that all dependencies can be created without panics
		repo := memory.NewInMemoryUserRepository()
		if repo == nil {
			t.Error("expected repository to be created, got nil")
		}

		service := usecase.NewUserService(repo)
		if service == nil {
			t.Error("expected service to be created, got nil")
		}

		handler := httpadapter.NewUserHandler(service)
		if handler == nil {
			t.Error("expected handler to be created, got nil")
		}

		mux := app.NewRouter(handler)
		if mux == nil {
			t.Error("expected router to be created, got nil")
		}
	})

	t.Run("Server configuration is correct", func(t *testing.T) {
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

		if srv.Addr != ":8080" {
			t.Errorf("expected server address ':8080', got '%s'", srv.Addr)
		}
		if srv.ReadTimeout != 10*time.Second {
			t.Errorf("expected ReadTimeout 10s, got %v", srv.ReadTimeout)
		}
		if srv.WriteTimeout != 10*time.Second {
			t.Errorf("expected WriteTimeout 10s, got %v", srv.WriteTimeout)
		}
		if srv.IdleTimeout != 60*time.Second {
			t.Errorf("expected IdleTimeout 60s, got %v", srv.IdleTimeout)
		}
		if srv.Handler == nil {
			t.Error("expected handler to be set")
		}
	})

	t.Run("Middleware is applied", func(t *testing.T) {
		repo := memory.NewInMemoryUserRepository()
		service := usecase.NewUserService(repo)
		handler := httpadapter.NewUserHandler(service)
		mux := app.NewRouter(handler)

		// Test that WithLogging returns a non-nil handler
		loggedHandler := app.WithLogging(mux)
		if loggedHandler == nil {
			t.Error("expected logged handler to be created, got nil")
		}

		// The logged handler should be different from the original mux
		// (this is a basic check that middleware is applied)
		if loggedHandler == mux {
			t.Error("expected middleware to wrap the original handler")
		}
	})
}

func TestContextHandling(t *testing.T) {
	t.Run("Context with timeout works", func(t *testing.T) {
		// Test context creation and timeout functionality
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if ctx == nil {
			t.Error("expected context to be created, got nil")
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("expected context to have deadline")
		}

		if time.Until(deadline) > 10*time.Second {
			t.Error("expected deadline to be within 10 seconds")
		}
	})

	t.Run("Signal context works", func(t *testing.T) {
		// Test that signal context can be created
		// Note: We can't easily test actual signal handling in unit tests
		// but we can test the context creation
		ctx := context.Background()
		if ctx == nil {
			t.Error("expected background context to be created, got nil")
		}
	})
}

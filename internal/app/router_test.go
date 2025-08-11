package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "cleanarch/internal/adapter/http"
	"cleanarch/internal/repository/memory"
	"cleanarch/internal/usecase"
)

func TestNewRouter(t *testing.T) {
	t.Run("Router is created successfully", func(t *testing.T) {
		// Create actual dependencies instead of empty service
		repo := memory.NewInMemoryUserRepository()
		service := usecase.NewUserService(repo)
		handler := httpadapter.NewUserHandler(service)

		router := NewRouter(handler)
		if router == nil {
			t.Error("expected router to be created, got nil")
		}
	})

	t.Run("Health check endpoint works", func(t *testing.T) {
		repo := memory.NewInMemoryUserRepository()
		service := usecase.NewUserService(repo)
		handler := httpadapter.NewUserHandler(service)
		router := NewRouter(handler)

		req := httptest.NewRequest("GET", "/healthz", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		expectedBody := "ok"
		if w.Body.String() != expectedBody {
			t.Errorf("expected body '%s', got '%s'", expectedBody, w.Body.String())
		}
	})

	t.Run("User endpoints are registered", func(t *testing.T) {
		repo := memory.NewInMemoryUserRepository()
		service := usecase.NewUserService(repo)
		handler := httpadapter.NewUserHandler(service)
		router := NewRouter(handler)

		// Test basic endpoints without path parameters
		testCases := []struct {
			method string
			path   string
		}{
			{"POST", "/api/v1/users"},
			{"GET", "/api/v1/users"},
		}

		for _, tc := range testCases {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// We expect the handler to be called (not 404)
			if w.Code == http.StatusNotFound {
				t.Errorf("endpoint %s %s not registered, got 404", tc.method, tc.path)
			}
		}

		// Note: Path parameter endpoints (GET /api/v1/users/{id}, etc.)
		// are tested in integration tests or through the actual HTTP handlers
		// since the Go 1.22 path parameter feature requires a real HTTP server context
	})

	t.Run("Non-existent endpoint returns 404", func(t *testing.T) {
		repo := memory.NewInMemoryUserRepository()
		service := usecase.NewUserService(repo)
		handler := httpadapter.NewUserHandler(service)
		router := NewRouter(handler)

		req := httptest.NewRequest("GET", "/non-existent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d for non-existent endpoint, got %d", http.StatusNotFound, w.Code)
		}
	})
}

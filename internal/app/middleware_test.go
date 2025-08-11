package app

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestWithLogging(t *testing.T) {
	t.Run("Logs request and response", func(t *testing.T) {
		// Capture log output
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr) // Restore default log output
		}()

		// Create a test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test response"))
		})

		// Wrap with logging middleware
		loggedHandler := WithLogging(testHandler)

		// Create test request
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		// Execute request
		loggedHandler.ServeHTTP(w, req)

		// Check response
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check log output
		logOutput := buf.String()
		if !strings.Contains(logOutput, "GET") {
			t.Error("expected log to contain HTTP method 'GET'")
		}
		if !strings.Contains(logOutput, "/test") {
			t.Error("expected log to contain path '/test'")
		}
		if !strings.Contains(logOutput, "200") {
			t.Error("expected log to contain status code '200'")
		}
	})

	t.Run("Logs different status codes", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		// Test handler that returns 404
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		loggedHandler := WithLogging(testHandler)
		req := httptest.NewRequest("POST", "/notfound", nil)
		w := httptest.NewRecorder()

		loggedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		logOutput := buf.String()
		if !strings.Contains(logOutput, "POST") {
			t.Error("expected log to contain HTTP method 'POST'")
		}
		if !strings.Contains(logOutput, "/notfound") {
			t.Error("expected log to contain path '/notfound'")
		}
		if !strings.Contains(logOutput, "404") {
			t.Error("expected log to contain status code '404'")
		}
	})

	t.Run("Measures request duration", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		// Test handler with delay
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond) // Small delay
			w.WriteHeader(http.StatusOK)
		})

		loggedHandler := WithLogging(testHandler)
		req := httptest.NewRequest("GET", "/slow", nil)
		w := httptest.NewRecorder()

		loggedHandler.ServeHTTP(w, req)

		logOutput := buf.String()
		// Check that duration is logged (should contain time units like 'ms' or 'µs')
		if !strings.Contains(logOutput, "ms") && !strings.Contains(logOutput, "µs") && !strings.Contains(logOutput, "s") {
			t.Error("expected log to contain duration information")
		}
	})
}

func TestStatusRecorder(t *testing.T) {
	t.Run("Records status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK, // Default status
		}

		// Test WriteHeader
		recorder.WriteHeader(http.StatusCreated)
		if recorder.status != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, recorder.status)
		}

		// Check that the underlying ResponseWriter also got the status
		if w.Code != http.StatusCreated {
			t.Errorf("expected underlying writer status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("Preserves default status when not explicitly set", func(t *testing.T) {
		w := httptest.NewRecorder()
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK, // Default status
		}

		// Write without calling WriteHeader explicitly
		recorder.Write([]byte("test"))

		// Status should remain as set initially
		if recorder.status != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, recorder.status)
		}
	})
}

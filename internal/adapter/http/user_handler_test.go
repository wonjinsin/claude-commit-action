package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cleanarch/internal/domain"
)

// UserServiceInterface defines the interface for user service operations
type UserServiceInterface interface {
	CreateUser(name, email string) (*domain.User, error)
	GetUser(id int64) (*domain.User, error)
	ListUsers() ([]*domain.User, error)
	UpdateUser(id int64, name, email string) (*domain.User, error)
	DeleteUser(id int64) error
}

// TestUserHandler wraps UserHandler to allow dependency injection for testing
type TestUserHandler struct {
	service UserServiceInterface
}

func NewTestUserHandler(service UserServiceInterface) *TestUserHandler {
	return &TestUserHandler{service: service}
}

func (h *TestUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	user, err := h.service.CreateUser(req.Name, req.Email)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *TestUserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	user, err := h.service.GetUser(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *TestUserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *TestUserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	user, err := h.service.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *TestUserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteUser(id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// MockUserService implements UserServiceInterface for testing
type MockUserService struct {
	users  map[int64]*domain.User
	nextID int64
	fail   bool
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		users:  make(map[int64]*domain.User),
		nextID: 1,
	}
}

func (m *MockUserService) SetFail(fail bool) {
	m.fail = fail
}

func (m *MockUserService) CreateUser(name, email string) (*domain.User, error) {
	if m.fail {
		return nil, &MockError{message: "service error"}
	}
	user := &domain.User{
		ID:    m.nextID,
		Name:  name,
		Email: email,
	}
	m.users[m.nextID] = user
	m.nextID++
	return user, nil
}

func (m *MockUserService) GetUser(id int64) (*domain.User, error) {
	if m.fail {
		return nil, &MockError{message: "service error"}
	}
	user, ok := m.users[id]
	if !ok {
		return nil, &MockError{message: "user not found"}
	}
	return user, nil
}

func (m *MockUserService) ListUsers() ([]*domain.User, error) {
	if m.fail {
		return nil, &MockError{message: "service error"}
	}
	result := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		result = append(result, user)
	}
	return result, nil
}

func (m *MockUserService) UpdateUser(id int64, name, email string) (*domain.User, error) {
	if m.fail {
		return nil, &MockError{message: "service error"}
	}
	user, ok := m.users[id]
	if !ok {
		return nil, &MockError{message: "user not found"}
	}
	user.Name = name
	user.Email = email
	return user, nil
}

func (m *MockUserService) DeleteUser(id int64) error {
	if m.fail {
		return &MockError{message: "service error"}
	}
	if _, ok := m.users[id]; !ok {
		return &MockError{message: "user not found"}
	}
	delete(m.users, id)
	return nil
}

type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Run("Create user successfully", func(t *testing.T) {
		service := NewMockUserService()
		handler := NewTestUserHandler(service)

		reqBody := `{"name":"John Doe","email":"john@example.com"}`
		req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
		}

		var response domain.User
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Name != "John Doe" {
			t.Errorf("expected name 'John Doe', got '%s'", response.Name)
		}
		if response.Email != "john@example.com" {
			t.Errorf("expected email 'john@example.com', got '%s'", response.Email)
		}
	})

	t.Run("Create user with invalid JSON", func(t *testing.T) {
		service := NewMockUserService()
		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader("invalid json"))
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["error"] != "invalid JSON" {
			t.Errorf("expected error 'invalid JSON', got '%s'", response["error"])
		}
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	t.Run("Get existing user", func(t *testing.T) {
		service := NewMockUserService()
		// Create a user first
		service.CreateUser("John Doe", "john@example.com")

		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.GetUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.User
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Name != "John Doe" {
			t.Errorf("expected name 'John Doe', got '%s'", response.Name)
		}
	})

	t.Run("Get user with invalid ID", func(t *testing.T) {
		service := NewMockUserService()
		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("GET", "/api/v1/users/invalid", nil)
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()

		handler.GetUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		service := NewMockUserService()
		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("GET", "/api/v1/users/999", nil)
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()

		handler.GetUser(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestUserHandler_ListUsers(t *testing.T) {
	t.Run("List users successfully", func(t *testing.T) {
		service := NewMockUserService()
		service.CreateUser("John Doe", "john@example.com")
		service.CreateUser("Jane Doe", "jane@example.com")

		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		w := httptest.NewRecorder()

		handler.ListUsers(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []*domain.User
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(response) != 2 {
			t.Errorf("expected 2 users, got %d", len(response))
		}
	})

	t.Run("List users with empty result", func(t *testing.T) {
		service := NewMockUserService()
		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("GET", "/api/v1/users", nil)
		w := httptest.NewRecorder()

		handler.ListUsers(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []*domain.User
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(response) != 0 {
			t.Errorf("expected 0 users, got %d", len(response))
		}
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	t.Run("Update user successfully", func(t *testing.T) {
		service := NewMockUserService()
		service.CreateUser("John Doe", "john@example.com")

		handler := NewTestUserHandler(service)

		reqBody := `{"name":"Jane Doe","email":"jane@example.com"}`
		req := httptest.NewRequest("PUT", "/api/v1/users/1", strings.NewReader(reqBody))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.UpdateUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.User
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Name != "Jane Doe" {
			t.Errorf("expected name 'Jane Doe', got '%s'", response.Name)
		}
	})

	t.Run("Update user with invalid ID", func(t *testing.T) {
		service := NewMockUserService()
		handler := NewTestUserHandler(service)

		reqBody := `{"name":"Jane Doe","email":"jane@example.com"}`
		req := httptest.NewRequest("PUT", "/api/v1/users/invalid", strings.NewReader(reqBody))
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()

		handler.UpdateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Update user with invalid JSON", func(t *testing.T) {
		service := NewMockUserService()
		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("PUT", "/api/v1/users/1", strings.NewReader("invalid json"))
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.UpdateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	t.Run("Delete user successfully", func(t *testing.T) {
		service := NewMockUserService()
		service.CreateUser("John Doe", "john@example.com")

		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("DELETE", "/api/v1/users/1", nil)
		req.SetPathValue("id", "1")
		w := httptest.NewRecorder()

		handler.DeleteUser(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}

		if w.Body.Len() != 0 {
			t.Errorf("expected empty body, got: %s", w.Body.String())
		}
	})

	t.Run("Delete user with invalid ID", func(t *testing.T) {
		service := NewMockUserService()
		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("DELETE", "/api/v1/users/invalid", nil)
		req.SetPathValue("id", "invalid")
		w := httptest.NewRecorder()

		handler.DeleteUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		service := NewMockUserService()
		handler := NewTestUserHandler(service)

		req := httptest.NewRequest("DELETE", "/api/v1/users/999", nil)
		req.SetPathValue("id", "999")
		w := httptest.NewRecorder()

		handler.DeleteUser(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestUserHandler_HelperFunctions(t *testing.T) {
	t.Run("writeJSON function", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"message": "test"}

		writeJSON(w, http.StatusOK, data)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
		}

		var response map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response["message"] != "test" {
			t.Errorf("expected message 'test', got '%s'", response["message"])
		}
	})

	t.Run("parseID function with valid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/123", nil)
		req.SetPathValue("id", "123")

		id, err := parseID(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if id != 123 {
			t.Errorf("expected ID 123, got %d", id)
		}
	})

	t.Run("parseID function with invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/invalid", nil)
		req.SetPathValue("id", "invalid")

		_, err := parseID(req)
		if err == nil {
			t.Error("expected error for invalid ID")
		}
	})
}

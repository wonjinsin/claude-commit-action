package app

import (
	httpadapter "cleanarch/internal/adapter/http"
	"net/http"
)

func NewRouter(userHandler *httpadapter.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/users", userHandler.CreateUser)
	mux.HandleFunc("GET /api/v1/users", userHandler.ListUsers)
	mux.HandleFunc("GET /api/v1/users/{id}", userHandler.GetUser)
	mux.HandleFunc("PUT /api/v1/users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /api/v1/users/{id}", userHandler.DeleteUser)

	// Healthcheck
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}

package handler

import (
	"encoding/json"
	"net/http"
	"sql-injection-subtle/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Search handles GET /users/search?q=... and GET /users?username=...
// It reads the search term from the request and passes it to the service.
// No database or SQL is referenced in this file.
func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Prefer "q" then "username" so analyzers see multiple sources
	q := r.URL.Query().Get("q")
	if q == "" {
		q = r.URL.Query().Get("username")
	}

	users, err := h.svc.Search(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"users": users})
}

// List handles GET /users?sort=...&order=...
// Sort and order are passed to the service and eventually into the query builder.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	sortColumn := r.URL.Query().Get("sort")
	if sortColumn == "" {
		sortColumn = "id"
	}
	sortDir := r.URL.Query().Get("order")

	users, err := h.svc.ListSorted(sortColumn, sortDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"users": users})
}

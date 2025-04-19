package tasks

import (
	"encoding/json"
	"net/http"

	"todo-app/internal/storage/postgresql"
)

type CreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateResponse struct {
	ID int `json:"id"`
}

func Create(store *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		id, err := store.CreateTask(r.Context(), req.Title, req.Description)
		if err != nil {
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateResponse{ID: id})
	}
}

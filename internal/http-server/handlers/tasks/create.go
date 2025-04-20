package tasks

import (
	"encoding/json"
	"net/http"
	"todo-app/internal/db/postgresql"
	"todo-app/internal/lib/api/response"
)

type CreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type CreateResponse struct {
	ID int `json:"id"`
}

func Create(store *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.Response{
				Status: response.StatusError,
				Error:  "invalid request body",
			})
			return
		}
		if req.Title == "" || req.Status == "" {
			response.WriteJSON(w, http.StatusBadRequest, response.Response{
				Status: response.StatusError,
				Error:  "title and status are required",
			})
			return
		}
		id, err := store.CreateTask(req.Title, req.Description, req.Status)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.Response{
				Status: response.StatusError,
				Error:  err.Error(),
			})
			return
		}
		response.WriteJSON(w, http.StatusCreated, CreateResponse{ID: id})
	}
}

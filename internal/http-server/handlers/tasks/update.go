package tasks

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"todo-app/internal/db/postgresql"
	"todo-app/internal/lib/api/response"
)

func Update(store *postgresql.Storage) http.HandlerFunc {
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
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.Response{
				Status: response.StatusError,
				Error:  err.Error(),
			})
			return
		}
		store.UpdateTask(id, req.Title, req.Description, req.Status)
		response.WriteJSON(w, http.StatusOK, "task with id - "+strconv.Itoa(id)+" has been updated")
	}
}

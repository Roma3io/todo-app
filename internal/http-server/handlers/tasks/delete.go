package tasks

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"todo-app/internal/db/postgresql"
	"todo-app/internal/lib/api/response"
)

func Delete(store *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.Response{
				Status: response.StatusError,
				Error:  err.Error(),
			})
			return
		}
		store.DeleteTask(id)
		response.WriteJSON(w, http.StatusOK, "deleted task"+strconv.Itoa(id))
	}
}

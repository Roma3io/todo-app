package tasks

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"todo-app/internal/db/postgresql"
	"todo-app/internal/lib/api/response"
)

func Get(st *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.Response{
				Status: response.StatusError,
				Error:  err.Error(),
			})
		}
		task, err := st.GetTask(id)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.Response{
				Status: response.StatusError,
				Error:  err.Error(),
			})
		}
		response.WriteJSON(w, http.StatusOK, task)
	}
}

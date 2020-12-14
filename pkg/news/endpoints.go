package news

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewNewsHandler(subRouter *mux.Router, newsSvc Service) http.Handler {
	subRouter.HandleFunc("/news", GetEventsHandler(newsSvc)).Methods(http.MethodGet)

	return subRouter
}

func GetEventsHandler(newsSvc Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()

		search := v.Get("search")

		news, err := newsSvc.GetNews(r.Context(), search)
		if err != nil {
			errorHandler(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, news)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, code int, err error) {
	respond(w, r, code, map[string]string{"error": err.Error()})
}

func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return
		}
	}
}
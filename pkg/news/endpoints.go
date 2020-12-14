package news

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewNewsHandler(subRouter *mux.Router, newsSvc Service, log *logrus.Logger) http.Handler {
	subRouter.HandleFunc("/news", GetEventsHandler(newsSvc, log)).Methods(http.MethodGet)

	return subRouter
}

func GetEventsHandler(newsSvc Service, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()

		search := v.Get("search")

		news, err := newsSvc.GetNews(r.Context(), search)
		if err != nil {
			log.Errorf("get news: %s", err)
			errorHandler(w, r, http.StatusUnprocessableEntity, err, log)
			return
		}

		respond(w, r, http.StatusOK, news, log)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, code int, err error, log *logrus.Logger) {
	respond(w, r, code, map[string]string{"error": err.Error()}, log)
}

func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}, log *logrus.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Errorf("respond json encode: %s", err)
			return
		}
	}
}
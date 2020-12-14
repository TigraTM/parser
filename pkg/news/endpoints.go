package news

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewNewsHandler(subRouter *mux.Router, newsSvc Service) http.Handler {
	return subRouter
}

func GetEventsHandler(newsSvc Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
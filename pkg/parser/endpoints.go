package parser

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func MakeParserHandler(subRouter *mux.Router, parserSvc Service) http.Handler {
	subRouter.HandleFunc("/parser", ParserHandler(parserSvc)).Methods(http.MethodPost)

	return subRouter
}

func ParserHandler(parserSvc Service) func(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Data Parser `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errorHandler(w, r, http.StatusBadRequest, err)
			return
		}

		if err := parserSvc.ParsingPage(r.Context(), req.Data); err != nil {
			errorHandler(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		respond(w, r, http.StatusOK, "Новости получены")
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
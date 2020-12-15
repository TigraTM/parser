package parser

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func MakeParserHandler(subRouter *mux.Router, parserSvc Service, log *logrus.Logger) http.Handler {
	subRouter.HandleFunc("/parser", ParserHandler(parserSvc, log)).Methods(http.MethodPost)

	return subRouter
}

func ParserHandler(parserSvc Service, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Data Parser `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Errorf("bad request: %s", err)
			errorHandler(w, r, http.StatusBadRequest, err, log)
			return
		}

		if err := parserSvc.ParsingPage(r.Context(), req.Data); err != nil {
			log.Errorf("parsing page: %s", err)
			errorHandler(w, r, http.StatusUnprocessableEntity, err, log)
			return
		}

		respond(w, r, http.StatusOK, nil, log)
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
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"parser/pkg/config"
	"parser/pkg/news"
	"parser/pkg/parser"
	"parser/pkg/storage"
)

func main() {
	cfg := config.New()

	newsRepo := storage.NewRepository(nil)
	newsSvc := news.NewService(newsRepo)
	parserSvc := parser.NewService(newsSvc)

	r := initRouter(newsSvc, parserSvc)
	srv := initServer(cfg, r)

	go func() {
		fmt.Println("server is running")
		if err := srv.ListenAndServe(); err != nil {
			fmt.Errorf("error listen and server: %w", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Errorf("error server shutdown: %w", err)
	}

	os.Exit(0)
}

func initRouter(newsSvc news.Service, parserSvc parser.Service) *mux.Router {
	r := mux.NewRouter()

	versionRout := r.PathPrefix("/v1").Subrouter()

	versionRout.HandleFunc("/hello", MakeHandler()).Methods(http.MethodGet)

	r.PathPrefix("/parser").Handler(parser.MakeParserHandler(versionRout, parserSvc))
	r.PathPrefix("/news").Handler(news.NewNewsHandler(versionRout, newsSvc))

	return r
}

func initServer(cfg *viper.Viper, r *mux.Router) *http.Server {
	return &http.Server{
		Addr:    cfg.GetString("LISTEN"),
		Handler: r,
	}
}

func MakeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]bool{"hello": true})
		if err != nil {
			log.Println(err)
		}
	}
}
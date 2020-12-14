package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"parser/migrations"
	"parser/pkg/config"
	"parser/pkg/news"
	"parser/pkg/parser"
	"parser/pkg/storage"
)

var log = logrus.New()

func main() {
	log.Out = os.Stdout
	log.SetFormatter(&logrus.JSONFormatter{})

	cfg := config.New()

	db, err := sqlx.Connect("postgres", cfg.GetString("DB"))
	if err != nil {
		log.Fatal(fmt.Errorf("db: %w", err))
	}

	db.MustExec(migrations.Schema)

	newsRepo := storage.NewRepository(db)
	newsSvc := news.NewService(newsRepo, log)
	parserSvc := parser.NewService(newsSvc, cfg, log)

	r := initRouter(newsSvc, parserSvc)
	srv := initServer(cfg, r)

	go func() {
		fmt.Println("server is running")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(fmt.Errorf("error listen and server: %w", err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Errorf("error server shutdown: %w", err))
	}

	os.Exit(0)
}

func initRouter(newsSvc news.Service, parserSvc parser.Service) *mux.Router {
	r := mux.NewRouter()

	versionRout := r.PathPrefix("/v1").Subrouter()

	versionRout.HandleFunc("/hello", MakeHandler()).Methods(http.MethodGet)

	r.PathPrefix("/parser").Handler(parser.MakeParserHandler(versionRout, parserSvc, log))

	r.PathPrefix("/news").Handler(news.NewNewsHandler(versionRout, newsSvc, log))

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

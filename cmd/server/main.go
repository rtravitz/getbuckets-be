package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/rtravitz/getbuckets-be/cmd/server/handler"
	"github.com/rtravitz/getbuckets-be/database"
)

type LogWriter struct {
	*log.Logger
}

func (lw LogWriter) Write(p []byte) (n int, err error) {
	lw.Print(string(p))
	return len(p), nil
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "localhost:5000"
	}

	return ":" + port
}

func run() error {
	// =========================================================================
	// Logging

	log := log.New(os.Stdout, "GETBUCKETS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:0.0.0.0"`
			Name       string `conf:"default:getbuckets"`
			DisableTLS bool   `conf:"default:false"`
		}
	}

	if err := conf.Parse(os.Args[1:], "GETBUCKETS", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("GETBUCKETS", &cfg)
			if err != nil {
				return err
			}
			fmt.Println(usage)
			return nil
		}
		return err
	}

	// =========================================================================
	// Start Database

	log.Println("main : Started : Initializing database support")

	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Printf("main : Database Stopping : %s", cfg.DB.Host)
		db.Close()
	}()

	// =========================================================================
	// Start API Service

	log.Println("main : Started : Initializing API support")

	api := API(log, db)

	port := getPort()
	log.Printf("Starting on port %s\n", port)

	return http.ListenAndServe(port, api)
}

// API ...
func API(log *log.Logger, db *sqlx.DB) http.Handler {
	r := mux.NewRouter()

	s := r.PathPrefix("/api/v0").Subrouter()
	s.HandleFunc("/buckets", handler.BucketsHandler(db)).Methods("GET")
	s.HandleFunc("/buckets", handler.SaveBucketHandler(db)).Methods("POST")
	s.HandleFunc("/buckets/{bucket_id}", handler.ShowBucketHandler(db)).Methods("GET")
	s.HandleFunc("/buckets/{bucket_id}/clean", handler.SaveCleanRatingHandler(db)).Methods("POST")
	s.HandleFunc("/buckets/{bucket_id}/lock", handler.SaveLockRatingHandler(db)).Methods("POST")

	lw := LogWriter{log}
	loggingRouter := handlers.LoggingHandler(lw, r)

	return cors.Default().Handler(loggingRouter)
}

func main() {
	log.Fatal(run())
}

package main

import (
	"log"
	"net/http"
	"time"
	"true_shop/internal/config"
	"true_shop/internal/db"
	"true_shop/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Database connection failed : %v", err)
	}
	log.Printf("Database listening")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", handlers.Listings(db))
	mux.HandleFunc("DELETE /listings/{id}",handlers.DeleteListing(db))
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Server listening %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed : %v", err)
	}

}

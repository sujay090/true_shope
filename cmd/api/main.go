package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
	"true_shop/internal/config"
	"true_shop/internal/db"
	"true_shop/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Database connection failed : %v", err)
	}

	log.Printf("Database listening")
	lh := handlers.NewListingHandler(db,logger)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)
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

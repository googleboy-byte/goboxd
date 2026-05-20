package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/handler"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	cfg, err := config.Load("languages.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Post("/run", handler.NewRunHandler(cfg))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

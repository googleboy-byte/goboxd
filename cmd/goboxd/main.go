package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/handler"
	"github.com/thesouldev/goboxd/internal/runner"
	"github.com/thesouldev/goboxd/internal/stats"
	"log/slog"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	configPath := flag.String("config", "languages.yaml", "path to languages.yaml")
	flag.Parse()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 1. Initialize stats
	s := stats.NewStats()

	// 2. Startup Cleanup
	runner.SweepOrphanedDirectories(os.TempDir(), 10*time.Minute)

	// 3. Startup Probes
	nsjailProbe := runner.ProbeNsjail()
	nsjailVer := nsjailProbe.Version
	if !nsjailProbe.OK {
		log.Printf("WARNING: nsjail probe failed: %s", nsjailProbe.Error)
	}

	langVers := make(map[string]string)
	for id, lang := range cfg.Languages {
		probe := runner.ProbeLanguage(lang)
		if probe.OK {
			langVers[id] = probe.Version
		} else {
			log.Printf("WARNING: language %s probe failed: %s", id, probe.Error)
			langVers[id] = "unknown"
		}
	}

	// 3. Handlers
	h := handler.NewHealthHandler(version, commit, nsjailVer, langVers, s, cfg)

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Get("/readyz", h.Readyz)
	r.Get("/info", h.Info)
	r.Post("/run", handler.NewRunHandler(cfg, s))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting %s (%s) on %s", version, commit, addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

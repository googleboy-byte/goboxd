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

func initCgroups() {
	if _, err := os.Stat("/sys/fs/cgroup/cgroup.subtree_control"); err != nil {
		slog.Warn("cgroup v2 subtree control not found, skipping initialization")
		return
	}

	// Move current process to a sub-cgroup to allow subtree control in root
	if err := os.MkdirAll("/sys/fs/cgroup/goboxd-node", 0755); err != nil {
		slog.Error("failed to create goboxd-node cgroup", "error", err)
		return
	}
	if err := os.WriteFile("/sys/fs/cgroup/goboxd-node/cgroup.procs", []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644); err != nil {
		slog.Error("failed to move to goboxd-node cgroup", "error", err)
		return
	}

	// Enable memory and pids controllers in root
	if err := os.WriteFile("/sys/fs/cgroup/cgroup.subtree_control", []byte("+memory +pids\n"), 0644); err != nil {
		slog.Error("failed to enable memory/pids in root", "error", err)
		return
	}

	// Prepare goboxd parent for nsjail with memory enabled
	if err := os.MkdirAll("/sys/fs/cgroup/goboxd", 0755); err != nil {
		slog.Error("failed to create goboxd parent cgroup", "error", err)
		return
	}
	if err := os.WriteFile("/sys/fs/cgroup/goboxd/cgroup.subtree_control", []byte("+memory +pids\n"), 0644); err != nil {
		slog.Error("failed to enable memory/pids in goboxd parent", "error", err)
	}
	slog.Info("cgroup initialization successful")
}

func main() {
	initCgroups()
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

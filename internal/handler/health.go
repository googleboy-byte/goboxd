package handler

import (
	"encoding/json"
	"net/http"
	"runtime"
	"syscall"
	"time"

	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/runner"
	"github.com/thesouldev/goboxd/internal/stats"
	"sync"
)

type HealthHandler struct {
	Version      string
	Commit       string
	NsjailPath   string
	NsjailVer    string
	LangVers     map[string]string
	Stats        *stats.Stats
	Config       *config.Config
	cache        probeCache
}

type probeCache struct {
	mu       sync.Mutex
	result   *ReadyzResponse
	cachedAt time.Time
	ttl      time.Duration
}

type NsjailStatus struct {
	OK      bool   `json:"ok"`
	Version string `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`
}

type LanguageStatus struct {
	OK      bool   `json:"ok"`
	Version string `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`
}

type NsjailInfo struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

type LanguageInfo struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Version          string        `json:"version"`
	DefaultRunLimits config.Limits `json:"default_run_limits"`
}

type LimitsInfo struct {
	MaxSourceBytes    int `json:"max_source_bytes"`
	MaxTests          int `json:"max_tests"`
	MaxConcurrentJobs int `json:"max_concurrent_jobs"`
}

type StatsInfo struct {
	InFlight          int64      `json:"in_flight_jobs"`
	QueueSize         int64      `json:"queue_size"`
	JobsTotal         int64      `json:"jobs_total"`
	JobsFailed        int64      `json:"jobs_failed_internal"`
	LastInternalErrAt *time.Time `json:"last_internal_error_at"`
	DiskFreeBytes     uint64     `json:"disk_free_bytes_jail_dir"`
}

type InfoResponse struct {
	BuildInfo struct {
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		GoVersion string `json:"go_version"`
	} `json:"build_info"`
	Nsjail    NsjailInfo     `json:"nsjail"`
	Languages []LanguageInfo `json:"languages"`
	Limits    LimitsInfo     `json:"limits"`
	Stats     StatsInfo      `json:"stats"`
}

type ReadyzResponse struct {
	Status    string                    `json:"status"`
	Nsjail    NsjailStatus              `json:"nsjail"`
	Languages map[string]LanguageStatus `json:"languages"`
}

func NewHealthHandler(version, commit, nsjailVer string, langVers map[string]string, s *stats.Stats, cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		Version:    version,
		Commit:     commit,
		NsjailPath: runner.NsjailPath,
		NsjailVer:  nsjailVer,
		LangVers:   langVers,
		Stats:      s,
		Config:     cfg,
		cache: probeCache{
			ttl: 30 * time.Second,
		},
	}
}

func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	h.cache.mu.Lock()
	defer h.cache.mu.Unlock()

	if h.cache.result != nil && time.Since(h.cache.cachedAt) < h.cache.ttl {
		h.writeReadyz(w, h.cache.result)
		return
	}

	nsjail := runner.ProbeNsjail()

	resp := ReadyzResponse{
		Languages: make(map[string]LanguageStatus),
	}
	resp.Nsjail = NsjailStatus{
		OK:      nsjail.OK,
		Version: nsjail.Version,
		Error:   nsjail.Error,
	}

	if !nsjail.OK {
		resp.Status = "degraded"
		h.cache.result = &resp
		h.cache.cachedAt = time.Now()
		h.writeReadyz(w, &resp)
		return
	}

	allOK := true
	for id, lang := range h.Config.Languages {
		probe := runner.ProbeLanguage(lang)
		resp.Languages[id] = LanguageStatus{
			OK:      probe.OK,
			Version: probe.Version,
			Error:   probe.Error,
		}
		if !probe.OK {
			allOK = false
		}
	}

	if allOK {
		resp.Status = "ok"
	} else {
		resp.Status = "degraded"
	}

	h.cache.result = &resp
	h.cache.cachedAt = time.Now()
	h.writeReadyz(w, &resp)
}

func (h *HealthHandler) writeReadyz(w http.ResponseWriter, resp *ReadyzResponse) {
	w.Header().Set("Content-Type", "application/json")
	if resp.Status == "ok" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *HealthHandler) sendDegraded(w http.ResponseWriter, resp ReadyzResponse) {
	resp.Status = "degraded"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(resp)
}

func (h *HealthHandler) Info(w http.ResponseWriter, r *http.Request) {
	resp := InfoResponse{}
	resp.BuildInfo.Version = h.Version
	resp.BuildInfo.Commit = h.Commit
	resp.BuildInfo.GoVersion = runtime.Version()

	resp.Nsjail = NsjailInfo{
		Path:    h.NsjailPath,
		Version: h.NsjailVer,
	}

	for id, lang := range h.Config.Languages {
		resp.Languages = append(resp.Languages, LanguageInfo{
			ID:               id,
			Name:             lang.Name,
			Version:          h.LangVers[id],
			DefaultRunLimits: lang.Run.Limits,
		})
	}

	resp.Limits = LimitsInfo{
		MaxSourceBytes:    262144,
		MaxTests:          50,
		MaxConcurrentJobs: h.Config.MaxConcurrentJobs,
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs("/tmp", &stat); err == nil {
		resp.Stats.DiskFreeBytes = stat.Bavail * uint64(stat.Bsize)
	}

	resp.Stats.InFlight = h.Stats.InFlight.Load()
	resp.Stats.QueueSize = h.Stats.QueueSize.Load()
	resp.Stats.JobsTotal = h.Stats.JobsTotal.Load()
	resp.Stats.JobsFailed = h.Stats.JobsFailedInternal.Load()
	resp.Stats.LastInternalErrAt = h.Stats.LastInternalErrAt.Load()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

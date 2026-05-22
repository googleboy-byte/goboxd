package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/runner"
	"github.com/thesouldev/goboxd/internal/stats"
	"github.com/thesouldev/goboxd/internal/validate"
	"time"
)

type ConfigOverride struct {
	Limits *config.Limits `json:"limits"`
	Flags  []string       `json:"flags"`
}

type Request struct {
	Language         string          `json:"language"`
	Source           string          `json:"source"`
	SourceFilename   string          `json:"source_filename"`
	ArtifactFilename string          `json:"artifact_filename"`
	Build            *ConfigOverride `json:"build"`
	Run              *ConfigOverride `json:"run"`
	Tests            []runner.TestCase `json:"tests"`
}

type BuildResult struct {
	Status     string `json:"status"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMs int64  `json:"duration_ms"`
}

type Response struct {
	Status string               `json:"status"`
	Build  *BuildResult         `json:"build,omitempty"`
	Tests  []runner.TestResult `json:"tests"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewRunHandler(cfg *config.Config, s *stats.Stats) http.HandlerFunc {
	sem := make(chan struct{}, cfg.MaxConcurrentJobs)
	return func(w http.ResponseWriter, r *http.Request) {
		s.JobsTotal.Add(1)

		// 1. Queueing
		select {
		case sem <- struct{}{}:
			// Acquired slot
		case <-r.Context().Done():
			return
		case <-time.After(time.Duration(cfg.QueueTimeoutS) * time.Second):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "queue_timeout",
					"message": "server is busy, try again later",
				},
			})
			return
		}

		s.InFlight.Add(1)
		defer func() {
			s.InFlight.Add(-1)
			<-sem
		}()

		r.Body = http.MaxBytesReader(w, r.Body, 256*1024)

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "invalid_json", err.Error())
			return
		}

		// 1. Language lookup
		lang, err := cfg.GetLanguage(req.Language)
		if err != nil {
			sendError(w, "unknown_language", err.Error())
			return
		}

		// 2. Validation
		if err := validate.ValidateRunRequest(req.Language, req.Source, len(req.Tests), 256*1024, 50); err != nil {
			sendError(w, "bad_request", err.Error())
			return
		}

		for _, tc := range req.Tests {
			if err := validate.ValidateTest(tc.Stdin, tc.ExpectedOutput, 64*1024, 64*1024); err != nil {
				sendError(w, "bad_request", err.Error())
				return
			}
		}

		if req.SourceFilename != "" {
			if err := validate.ValidateFilename(req.SourceFilename); err != nil {
				sendError(w, "invalid_filename", err.Error())
				return
			}
			lang.SourceFilename = req.SourceFilename
		} else if err := validate.ValidateFilename(lang.SourceFilename); err != nil {
			sendError(w, "invalid_filename", err.Error())
			return
		}

		if req.ArtifactFilename != "" {
			if err := validate.ValidateFilename(req.ArtifactFilename); err != nil {
				sendError(w, "invalid_filename", err.Error())
				return
			}
			lang.Artifact = req.ArtifactFilename
		}

		if req.Build != nil && lang.Build != nil {
			if err := validate.ValidateFlags(req.Build.Flags, lang.Build.FlagAllowlist); err != nil {
				sendError(w, "disallowed_flag", err.Error())
				return
			}
			// Deep copy Build to avoid mutating global config
			cp := *lang.Build
			lang.Build = &cp
			if req.Build.Limits != nil {
				l := req.Build.Limits
				if l.WallTimeS > 0 {
					lang.Build.Limits.WallTimeS = l.WallTimeS
				}
				if l.MemoryKB > 0 {
					lang.Build.Limits.MemoryKB = l.MemoryKB
				}
				if l.MaxProcesses > 0 {
					lang.Build.Limits.MaxProcesses = l.MaxProcesses
				}
			}
		}

		var runFlags []string
		if req.Run != nil {
			if err := validate.ValidateFlags(req.Run.Flags, lang.Run.FlagAllowlist); err != nil {
				sendError(w, "disallowed_flag", err.Error())
				return
			}
			runFlags = req.Run.Flags
			if req.Run.Limits != nil {
				l := req.Run.Limits
				if l.WallTimeS > 0 {
					lang.Run.Limits.WallTimeS = l.WallTimeS
				}
				if l.MemoryKB > 0 {
					lang.Run.Limits.MemoryKB = l.MemoryKB
				}
				if l.MaxProcesses > 0 {
					lang.Run.Limits.MaxProcesses = l.MaxProcesses
				}
			}
		}

		var buildFlags []string
		if req.Build != nil {
			buildFlags = req.Build.Flags
		}

		// 3. Execution
		// Mapping handler.Request to runner.RunRequest
		runReq := runner.RunRequest{
			Source:     req.Source,
			Tests:      req.Tests,
			BuildFlags: buildFlags,
			RunFlags:   runFlags,
		}

		runResult := runner.Run(lang, runReq)

		// 4. Response
		resp := Response{
			Status: runResult.Status,
			Tests:  runResult.TestResults,
		}

		if runResult.Status == "internal_error" {
			s.JobsFailedInternal.Add(1)
			now := time.Now()
			s.LastInternalErrAt.Store(&now)
		}

		if runResult.Build != nil {
			resp.Build = &BuildResult{
				Status:     runResult.Build.Status,
				Stdout:     runResult.Build.Stdout,
				Stderr:     runResult.Build.Stderr,
				DurationMs: runResult.Build.DurationMs,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func sendError(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	var resp ErrorResponse
	resp.Error.Code = code
	resp.Error.Message = message
	json.NewEncoder(w).Encode(resp)
}

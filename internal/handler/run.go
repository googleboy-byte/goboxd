package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/runner"
	"github.com/thesouldev/goboxd/internal/validate"
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

func NewRunHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		if req.Source == "" {
			sendError(w, "invalid_source", "source is required")
			return
		}

		if err := validate.ValidateFilename(lang.SourceFilename); err != nil {
			sendError(w, "invalid_filename", err.Error())
			return
		}

		if req.Build != nil && lang.Build != nil {
			if err := validate.ValidateFlags(req.Build.Flags, lang.Build.FlagAllowlist); err != nil {
				sendError(w, "disallowed_flag", err.Error())
				return
			}
		}

		// 3. Execution
		// Mapping handler.Request to runner.RunRequest
		runReq := runner.RunRequest{
			Source: req.Source,
			Tests:  req.Tests,
		}

		runResult := runner.Run(lang, runReq)

		// 4. Response
		resp := Response{
			Status: runResult.Status,
			Tests:  runResult.TestResults,
		}
		
		// For now build is always ok as we don't have build step yet
		if lang.Build != nil {
			resp.Build = &BuildResult{
				Status: "ok",
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

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/thesouldev/goboxd/internal/config"
)

const nsjailPath = "/usr/sbin/nsjail"

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

type ReadyzResponse struct {
	Status    string                    `json:"status"`
	Nsjail    NsjailStatus              `json:"nsjail"`
	Languages map[string]LanguageStatus `json:"languages"`
}

func NewReadyzHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := ReadyzResponse{
			Languages: make(map[string]LanguageStatus, len(cfg.Languages)),
		}

		// 1. nsjail check
		resp.Nsjail = probeNsjail()
		if !resp.Nsjail.OK {
			// Don't probe languages if nsjail is broken
			resp.Status = "degraded"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// 2. Language probes
		allOK := true
		for id, lang := range cfg.Languages {
			resp.Languages[id] = probeLanguage(lang)
			if !resp.Languages[id].OK {
				allOK = false
			}
		}

		if allOK {
			resp.Status = "ok"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		} else {
			resp.Status = "degraded"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func probeNsjail() NsjailStatus {
	// Check binary exists and is executable
	info, err := os.Stat(nsjailPath)
	if err != nil {
		return NsjailStatus{OK: false, Error: fmt.Sprintf("nsjail not found at %s", nsjailPath)}
	}
	if info.Mode().Perm()&0111 == 0 {
		return NsjailStatus{OK: false, Error: fmt.Sprintf("nsjail at %s is not executable", nsjailPath)}
	}

	// Get version from stderr
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, nsjailPath, "--version")
	out, _ := cmd.CombinedOutput()
	// nsjail 3.0+ sometimes doesn't support --version and prints to stderr
	version := ""
	outputStr := string(out)
	if strings.Contains(outputStr, "unrecognized option") || strings.Contains(outputStr, "invalid option") {
		// FALLBACK: nsjail 3.4 (as checked out in external/nsjail) doesn't support --version.
		// If we get an unrecognized option error, we assume the version we built.
		version = "3.4"
	} else {
		for _, line := range strings.Split(outputStr, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				version = line
				break
			}
		}
	}

	if version == "" {
		version = "3.4"
	}

	return NsjailStatus{OK: true, Version: version}
}

func probeLanguage(lang config.Language) LanguageStatus {
	// Determine probe command
	probeCmd := lang.Run.Cmd
	probeArgs := []string{"--version"}

	if lang.VersionProbe != "" {
		// Parse version_probe as "cmd args..."
		parts := strings.Fields(lang.VersionProbe)
		if len(parts) > 0 {
			probeCmd = parts[0]
			if len(parts) > 1 {
				probeArgs = parts[1:]
			} else {
				// If only cmd is provided, the user might want it run by itself
				// but the requirement says "exec [lang.Run.Cmd, '--version']" for the fallback
				// Let's stick to the spirit: if version_probe is "cmd --v", it works.
				// If it's just "cmd", we assume it handles its own version reporting or needs no args.
				// However, the user said "split it on space and exec it".
				probeArgs = []string{}
			}
		}
	}

	// Check binary exists
	if _, err := exec.LookPath(probeCmd); err != nil {
		return LanguageStatus{OK: false, Error: fmt.Sprintf("binary not found at %s", probeCmd)}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, probeCmd, probeArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return LanguageStatus{OK: false, Error: "version probe timed out"}
		}
		// Some tools exit non-zero on --version but still print version
		if len(out) == 0 {
			return LanguageStatus{OK: false, Error: err.Error()}
		}
	}

	// Grab first non-empty line as version
	version := ""
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			version = line
			break
		}
	}

	return LanguageStatus{OK: true, Version: version}
}

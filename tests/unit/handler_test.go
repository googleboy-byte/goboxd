package unit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/handler"
	"github.com/thesouldev/goboxd/internal/stats"
)

func testConfig() *config.Config {
	return &config.Config{
		Languages: map[string]config.Language{
			"py3": {
				ID:             "py3",
				Name:           "Python 3",
				SourceFilename: "solution.py",
				Run: config.RunConfig{
					Cmd:    "/usr/bin/python3",
					Args:   []string{"{{source}}"},
					Limits: config.Limits{WallTimeS: 9, MemoryKB: 102400, MaxProcesses: 100},
				},
			},
			"cpp": {
				ID:             "cpp",
				Name:           "C++",
				SourceFilename: "solution.cpp",
				Artifact:       "solution",
				Build: &config.BuildConfig{
					Cmd:           "/usr/bin/g++",
					Args:          []string{"{{flags}}", "-o", "{{artifact}}", "{{source}}"},
					Limits:        config.Limits{WallTimeS: 3, MemoryKB: 1048576, MaxProcesses: 100},
					FlagAllowlist: []string{"-O0", "-O1", "-O2", "-O3", "-Wall", "-std=*"},
				},
				Run: config.RunConfig{
					Cmd:    "./{{artifact}}",
					Limits: config.Limits{WallTimeS: 3, MemoryKB: 524288, MaxProcesses: 64},
				},
			},
		},
		MaxConcurrentJobs: 4,
		QueueTimeoutS:     30,
	}
}

func TestRunHandler_UnknownLanguage(t *testing.T) {
	cfg := testConfig()
	s := stats.NewStats()
	h := handler.NewRunHandler(cfg, s)

	body := `{"language":"fortran","source":"print *,'hi'","tests":[{"stdin":"","expected_stdout":"hi"}]}`
	req := httptest.NewRequest(http.MethodPost, "/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "unknown_language") {
		t.Errorf("expected error code unknown_language, got %s", w.Body.String())
	}
}

func TestRunHandler_MalformedJSON(t *testing.T) {
	cfg := testConfig()
	s := stats.NewStats()
	h := handler.NewRunHandler(cfg, s)

	body := `{not valid json`
	req := httptest.NewRequest(http.MethodPost, "/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "invalid_json") {
		t.Errorf("expected error code invalid_json, got %s", w.Body.String())
	}
}

func TestRunHandler_DisallowedFlag(t *testing.T) {
	cfg := testConfig()
	s := stats.NewStats()
	h := handler.NewRunHandler(cfg, s)

	body := `{
		"language":"cpp",
		"source":"#include <iostream>\nint main(){return 0;}",
		"tests":[{"stdin":"","expected_stdout":""}],
		"build":{"flags":["-fplugin=evil.so"]}
	}`
	req := httptest.NewRequest(http.MethodPost, "/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "disallowed_flag") {
		t.Errorf("expected error code disallowed_flag, got %s", w.Body.String())
	}
}

func TestRunHandler_EmptySource(t *testing.T) {
	cfg := testConfig()
	s := stats.NewStats()
	h := handler.NewRunHandler(cfg, s)

	body := `{"language":"py3","source":"","tests":[{"stdin":"","expected_stdout":""}]}`
	req := httptest.NewRequest(http.MethodPost, "/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "bad_request") {
		t.Errorf("expected error code bad_request, got %s", w.Body.String())
	}
}

func TestRunHandler_TooManyTests(t *testing.T) {
	cfg := testConfig()
	s := stats.NewStats()
	h := handler.NewRunHandler(cfg, s)

	// Build JSON with 51 tests
	tests := make([]string, 51)
	for i := range tests {
		tests[i] = `{"stdin":"","expected_stdout":""}`
	}
	body := `{"language":"py3","source":"print()","tests":[` + strings.Join(tests, ",") + `]}`
	req := httptest.NewRequest(http.MethodPost, "/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "bad_request") {
		t.Errorf("expected error code bad_request, got %s", w.Body.String())
	}
}

func TestRunHandler_NoTests(t *testing.T) {
	cfg := testConfig()
	s := stats.NewStats()
	h := handler.NewRunHandler(cfg, s)

	body := `{"language":"py3","source":"print()","tests":[]}`
	req := httptest.NewRequest(http.MethodPost, "/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRunHandler_OversizeBody(t *testing.T) {
	cfg := testConfig()
	s := stats.NewStats()
	h := handler.NewRunHandler(cfg, s)

	// 300KiB source — exceeds 256KiB MaxBytesReader
	bigSource := strings.Repeat("x", 300*1024)
	body := `{"language":"py3","source":"` + bigSource + `","tests":[{"stdin":"","expected_stdout":""}]}`
	req := httptest.NewRequest(http.MethodPost, "/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRunHandler_QueueTimeout(t *testing.T) {
	cfg := testConfig()
	// Set MaxConcurrentJobs to 0 to make the semaphore channel unbuffered/blocking.
	// Set QueueTimeoutS to a very small value for the test.
	cfg.MaxConcurrentJobs = 0
	cfg.QueueTimeoutS = 1

	s := stats.NewStats()
	h := handler.NewRunHandler(cfg, s)

	body := `{"language":"py3","source":"print()","tests":[{"stdin":"","expected_stdout":""}]}`
	req := httptest.NewRequest(http.MethodPost, "/run", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// This should timeout after 1s because MaxConcurrentJobs=0 
	// mean there is no space in the unbuffered channel.
	h.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "queue_timeout") {
		t.Errorf("expected error code queue_timeout, got %s", w.Body.String())
	}
}

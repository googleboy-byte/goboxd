package runner

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/thesouldev/goboxd/internal/config"
)

type TestCase struct {
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_stdout"`
}

type RunRequest struct {
	Source string     `json:"source"`
	Tests  []TestCase `json:"tests"`
}

type TestResult struct {
	Status       string `json:"status"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	DurationMs   int64  `json:"duration_ms"`
	MemoryPeakKB int64  `json:"memory_peak_kb"`
}

type RunResult struct {
	Status      string       `json:"status"` // accepted, rejected
	TestResults []TestResult `json:"test_results"`
}

func Run(lang config.Language, req RunRequest) RunResult {
	workdir, err := os.MkdirTemp("", "goboxd-*")
	if err != nil {
		return RunResult{Status: "internal_error"}
	}
	defer os.RemoveAll(workdir)

	sourcePath := filepath.Join(workdir, lang.SourceFilename)
	if err := os.WriteFile(sourcePath, []byte(req.Source), 0644); err != nil {
		return RunResult{Status: "internal_error"}
	}

	results := make([]TestResult, 0, len(req.Tests))
	overallStatus := "accepted"

	vars := map[string]string{
		"source":   "/sandbox/" + lang.SourceFilename,
		"artifact": "/sandbox/" + lang.Artifact,
	}
	runArgs := resolveArgs(lang.Run.Args, vars)

	for _, tc := range req.Tests {
		res := runTestCase(lang, workdir, runArgs, tc)
		results = append(results, res)
		if res.Status != "accepted" && overallStatus == "accepted" {
			overallStatus = res.Status // Set first failing status
		}
	}

	return RunResult{
		Status:      overallStatus,
		TestResults: results,
	}
}

func runTestCase(lang config.Language, workdir string, runArgs []string, tc TestCase) TestResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(lang.Run.Limits.WallTimeS+1)*time.Second)
	defer cancel()

	nsjailArgs := buildNsjailArgs(lang, workdir, lang.Run.Cmd, runArgs)
	cmd := exec.CommandContext(ctx, nsjailPath, nsjailArgs...)
	
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return TestResult{Status: "internal_error"}
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return TestResult{Status: "internal_error"}
	}

	cmd.Stdin = bytes.NewBufferString(tc.Stdin)

	start := time.Now()
	if err := cmd.Start(); err != nil {
		return TestResult{Status: "runtime_error"}
	}

	var stdout, stderr bytes.Buffer
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go func() {
		io.Copy(&stdout, io.LimitReader(stdoutPipe, 1024*64))
		stdoutDone <- struct{}{}
	}()
	go func() {
		io.Copy(&stderr, io.LimitReader(stderrPipe, 1024*64))
		stderrDone <- struct{}{}
	}()

	<-stdoutDone
	<-stderrDone
	waitErr := cmd.Wait()
	duration := time.Since(start).Milliseconds()

	if ctx.Err() == context.DeadlineExceeded {
		return TestResult{Status: "time_exceeded", DurationMs: duration}
	}

	status := "accepted"
	if waitErr != nil {
		status = "runtime_error"
	} else {
		actual := stdout.String()
		expected := tc.ExpectedOutput

		if actual == expected {
			status = "accepted"
		} else if strings.TrimSpace(actual) == strings.TrimSpace(expected) {
			status = "output_whitespace_mismatch"
		} else {
			status = "wrong_output"
		}
	}

	return TestResult{
		Status:     status,
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		DurationMs: duration,
	}
}

func resolveArgs(args []string, vars map[string]string) []string {
	out := make([]string, len(args))
	for i, a := range args {
		for k, v := range vars {
			a = strings.ReplaceAll(a, "{{"+k+"}}", v)
		}
		out[i] = a
	}
	return out
}

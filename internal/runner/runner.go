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
	Status string `json:"status"` // accepted, wrong_answer, time_limit_exceeded, runtime_error
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
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

	for _, tc := range req.Tests {
		res := runTestCase(lang, sourcePath, tc)
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

func runTestCase(lang config.Language, sourcePath string, tc TestCase) TestResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(lang.Run.Limits.WallTimeS)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/usr/bin/python3", sourcePath)
	
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return TestResult{Status: "runtime_error", Error: err.Error()}
	}

	cmd.Stdin = bytes.NewBufferString(tc.Stdin)

	if err := cmd.Start(); err != nil {
		return TestResult{Status: "runtime_error", Error: err.Error()}
	}

	var stdout bytes.Buffer
	_, err = io.Copy(&stdout, io.LimitReader(stdoutPipe, 1024*64))
	
	// We still need to call Wait to clean up the process
	waitErr := cmd.Wait()

	if ctx.Err() == context.DeadlineExceeded {
		return TestResult{Status: "time_exceeded"}
	}

	if waitErr != nil {
		return TestResult{Status: "runtime_error", Error: waitErr.Error()}
	}

	actualOutput := strings.TrimSpace(stdout.String())
	expectedOutput := strings.TrimSpace(tc.ExpectedOutput)

	if actualOutput != expectedOutput {
		return TestResult{Status: "wrong_output", Output: actualOutput}
	}

	return TestResult{Status: "accepted", Output: actualOutput}
}

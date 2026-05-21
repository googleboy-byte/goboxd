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
	Source     string      `json:"source"`
	Tests      []TestCase  `json:"tests"`
	BuildFlags []string    `json:"build_flags"`
}

type BuildResult struct {
	Status     string `json:"status"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMs int64  `json:"duration_ms"`
}

type TestResult struct {
	Status       string `json:"status"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	DurationMs   int64  `json:"duration_ms"`
	MemoryPeakKB int64  `json:"memory_peak_kb"`
}

type RunResult struct {
	Status      string       `json:"status"` // accepted, build_failed, rejected
	Build       *BuildResult `json:"build"`
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

	var buildRes *BuildResult
	if lang.Build != nil {
		res := buildArtifact(lang, workdir, req.BuildFlags)
		buildRes = &res
		if res.Status != "ok" {
			notExecuted := make([]TestResult, len(req.Tests))
			for i := range notExecuted {
				notExecuted[i] = TestResult{Status: "not_executed"}
			}
			return RunResult{
				Status:      "build_failed",
				Build:       buildRes,
				TestResults: notExecuted,
			}
		}
	}

	vars := map[string]string{
		"source":   "/sandbox/" + lang.SourceFilename,
		"artifact": lang.Artifact,
	}
	runCmd := ResolveString(lang.Run.Cmd, vars)
	runArgs := ResolveArgs(lang.Run.Args, vars)

	for _, tc := range req.Tests {
		res := runTestCase(lang, workdir, runCmd, runArgs, tc)
		results = append(results, res)
		if res.Status != "accepted" && overallStatus == "accepted" {
			overallStatus = res.Status // Set first failing status
		}
	}

	return RunResult{
		Status:      overallStatus,
		Build:       buildRes,
		TestResults: results,
	}
}

func buildArtifact(lang config.Language, workdir string, extraFlags []string) BuildResult {
	vars := map[string]string{
		"source":   "/sandbox/" + lang.SourceFilename,
		"artifact": "/sandbox/" + lang.Artifact,
		"flags":    "",
	}
	resolved := ResolveArgs(lang.Build.Args, vars)
	finalArgs := make([]string, 0, len(resolved)+len(extraFlags))
	finalArgs = append(finalArgs, extraFlags...)
	for _, arg := range resolved {
		if arg != "" {
			finalArgs = append(finalArgs, arg)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(lang.Build.Limits.WallTimeS+1)*time.Second)
	defer cancel()

	buildCmd := ResolveString(lang.Build.Cmd, vars)
	nsjailArgs := buildNsjailArgsBuild(lang, workdir, buildCmd, finalArgs)
	cmd := exec.CommandContext(ctx, nsjailPath, nsjailArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start).Milliseconds()

	status := "ok"
	if err != nil {
		status = "failed"
	}

	return BuildResult{
		Status:     status,
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		DurationMs: duration,
	}
}

func runTestCase(lang config.Language, workdir string, runCmd string, runArgs []string, tc TestCase) TestResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(lang.Run.Limits.WallTimeS+1)*time.Second)
	defer cancel()

	nsjailArgs := buildNsjailArgs(lang, workdir, runCmd, runArgs)
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

	const limit = 1024 * 64
	const marker = "\n[TRUNCATED]\n"

	go func() {
		n, _ := io.Copy(&stdout, io.LimitReader(stdoutPipe, limit))
		if n >= limit {
			stdout.WriteString(marker)
		}
		stdoutDone <- struct{}{}
	}()
	go func() {
		n, _ := io.Copy(&stderr, io.LimitReader(stderrPipe, limit))
		if n >= limit {
			stderr.WriteString(marker)
		}
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

func ResolveArgs(args []string, vars map[string]string) []string {
	out := make([]string, len(args))
	for i, a := range args {
		out[i] = ResolveString(a, vars)
	}
	return out
}

func ResolveString(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{{"+k+"}}", v)
	}
	return s
}

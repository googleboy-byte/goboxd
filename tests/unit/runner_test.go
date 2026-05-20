package unit

import (
	"testing"

	"github.com/thesouldev/goboxd/internal/config"
	"github.com/thesouldev/goboxd/internal/runner"
)

func TestRunHelloWorld(t *testing.T) {
	lang := config.Language{
		ID:             "py3",
		Name:           "Python 3",
		SourceFilename: "solution.py",
		Run: config.RunConfig{
			Cmd: "/usr/bin/python3",
			Limits: config.Limits{
				WallTimeS: 5,
			},
		},
	}

	req := runner.RunRequest{
		Source: "print('hello world')",
		Tests: []runner.TestCase{
			{
				Stdin:          "",
				ExpectedOutput: "hello world",
			},
		},
	}

	res := runner.Run(lang, req)

	if res.Status != "accepted" {
		t.Errorf("expected status accepted, got %s", res.Status)
		for _, tr := range res.TestResults {
			t.Logf("Test Status: %s, Output: %q, Error: %q", tr.Status, tr.Output, tr.Error)
		}
	}

	if len(res.TestResults) != 1 {
		t.Fatalf("expected 1 test result, got %d", len(res.TestResults))
	}

	if res.TestResults[0].Status != "accepted" {
		t.Errorf("expected test result status accepted, got %s", res.TestResults[0].Status)
	}
}

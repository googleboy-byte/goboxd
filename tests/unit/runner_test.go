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

	t.Run("Accepted on literal match", func(t *testing.T) {
		req := runner.RunRequest{
			Source: "print('hello world', end='')",
			Tests: []runner.TestCase{
				{
					Stdin:          "",
					ExpectedOutput: "hello world", // literal match
				},
			},
		}
		res := runner.Run(lang, req)
		if res.Status != "accepted" {
			t.Errorf("expected status accepted, got %s", res.Status)
		}
	})

	t.Run("Output whitespace mismatch", func(t *testing.T) {
		req := runner.RunRequest{
			Source: "print('hello world')",
			Tests: []runner.TestCase{
				{
					Stdin:          "",
					ExpectedOutput: "hello world", // mismatch because print adds \n
				},
			},
		}
		res := runner.Run(lang, req)
		if res.Status != "output_whitespace_mismatch" {
			t.Errorf("expected status output_whitespace_mismatch, got %s", res.Status)
		}
	})
}

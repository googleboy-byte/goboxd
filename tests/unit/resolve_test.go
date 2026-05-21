package unit

import (
	"testing"

	"github.com/thesouldev/goboxd/internal/runner"
)

func TestResolveString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		vars     map[string]string
		expected string
	}{
		{
			"Single placeholder",
			"/sandbox/{{source}}",
			map[string]string{"source": "solution.py"},
			"/sandbox/solution.py",
		},
		{
			"Multiple placeholders",
			"{{artifact}} from {{source}}",
			map[string]string{"artifact": "solution", "source": "solution.cpp"},
			"solution from solution.cpp",
		},
		{
			"No placeholders",
			"/usr/bin/python3",
			map[string]string{"source": "solution.py"},
			"/usr/bin/python3",
		},
		{
			"Unknown placeholder left as-is",
			"./{{unknown}}",
			map[string]string{"source": "solution.py"},
			"./{{unknown}}",
		},
		{
			"Empty string",
			"",
			map[string]string{"source": "solution.py"},
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.ResolveString(tt.input, tt.vars)
			if got != tt.expected {
				t.Errorf("ResolveString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestResolveArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		vars     map[string]string
		expected []string
	}{
		{
			"Resolves all args",
			[]string{"{{source}}", "-o", "{{artifact}}"},
			map[string]string{"source": "solution.cpp", "artifact": "solution"},
			[]string{"solution.cpp", "-o", "solution"},
		},
		{
			"Empty args",
			[]string{},
			map[string]string{"source": "solution.py"},
			[]string{},
		},
		{
			"No placeholders in args",
			[]string{"-Wall", "-Wextra"},
			map[string]string{"source": "solution.cpp"},
			[]string{"-Wall", "-Wextra"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runner.ResolveArgs(tt.args, tt.vars)
			if len(got) != len(tt.expected) {
				t.Fatalf("ResolveArgs length = %d, want %d", len(got), len(tt.expected))
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("ResolveArgs[%d] = %q, want %q", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

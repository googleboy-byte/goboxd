package validate

import (
	"testing"
)

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "solution.cpp", false},
		{"valid simple", "main", false},
		{"empty", "", true},
		{"too long", "a" + string(make([]byte, 65)), true},
		{"path separator slash", "dir/file", true},
		{"path separator backslash", "dir\\file", true},
		{"current dir", ".", true},
		{"parent dir", "..", true},
		{"traversal", "../etc/passwd", true},
		{"leading dot", ".hidden", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateFilename(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFlags(t *testing.T) {
	allowlist := []string{"-O0", "-O1", "-O2", "-O3", "-Wall", "-std=*"}
	tests := []struct {
		name      string
		requested []string
		wantErr   bool
	}{
		{"allowed exact", []string{"-O2", "-Wall"}, false},
		{"allowed glob", []string{"-std=c++17"}, false},
		{"disallowed", []string{"-fplugin=evil.so"}, true},
		{"mixed", []string{"-O2", "-fplugin=evil.so"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateFlags(tt.requested, allowlist); (err != nil) != tt.wantErr {
				t.Errorf("ValidateFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRunRequest(t *testing.T) {
	tests := []struct {
		name        string
		langId      string
		source      string
		testCount   int
		sourceLimit int
		testLimit   int
		wantErr     bool
	}{
		{"valid", "py3", "print(1)", 1, 256 * 1024, 50, false},
		{"empty lang", "", "print(1)", 1, 256 * 1024, 50, true},
		{"empty source", "py3", "", 1, 256 * 1024, 50, true},
		{"oversize source", "py3", "x", 1, 0, 50, true},
		{"zero tests", "py3", "print(1)", 0, 256 * 1024, 50, true},
		{"too many tests", "py3", "print(1)", 51, 256 * 1024, 50, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateRunRequest(tt.langId, tt.source, tt.testCount, tt.sourceLimit, tt.testLimit); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRunRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTest(t *testing.T) {
	tests := []struct {
		name          string
		stdin         string
		expected      string
		stdinLimit    int
		expectedLimit int
		wantErr       bool
	}{
		{"valid", "input", "output", 64 * 1024, 64 * 1024, false},
		{"oversize stdin", "x", "output", 0, 64 * 1024, true},
		{"oversize expected", "input", "x", 64 * 1024, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateTest(tt.stdin, tt.expected, tt.stdinLimit, tt.expectedLimit); (err != nil) != tt.wantErr {
				t.Errorf("ValidateTest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package unit

import (
	"errors"
	"testing"

	"github.com/thesouldev/goboxd/internal/validate"
)

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"Happy path", "solution.py", nil},
		{"Single component", "main", nil},
		{"Empty", "", validate.ErrInvalidFilename},
		{"Too long", "this_is_a_very_long_filename_that_exceeds_the_sixty_four_character_limit_defined_in_the_validation_logic.py", validate.ErrInvalidFilename},
		{"Path traversal ..", "../../etc/passwd", validate.ErrInvalidFilename},
		{"Path separator /", "dir/file.go", validate.ErrInvalidFilename},
		{"Reserved .", ".", validate.ErrInvalidFilename},
		{"Reserved ..", "..", validate.ErrInvalidFilename},
		{"Hidden file", ".git", validate.ErrInvalidFilename},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.ValidateFilename(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidateFlags(t *testing.T) {
	allowlist := []string{"-std=*", "-O[0-3]", "-Wall"}
	tests := []struct {
		name    string
		input   []string
		wantErr error
	}{
		{"Empty requested", []string{}, nil},
		{"All allowed", []string{"-std=c++17", "-O2", "-Wall"}, nil},
		{"Not allowed", []string{"-fplugin=evil.so"}, validate.ErrInvalidFlag},
		{"Mixed", []string{"-Wall", "-fplugin=evil.so"}, validate.ErrInvalidFlag},
		{"Empty allowlist", []string{"-Wall"}, validate.ErrInvalidFlag},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := allowlist
			if tt.name == "Empty allowlist" {
				al = []string{}
			}
			err := validate.ValidateFlags(tt.input, al)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidateRunRequest(t *testing.T) {
	tests := []struct {
		name      string
		langId    string
		source    string
		testCount int
		wantErr   error
	}{
		{"Happy path", "py3", "print('hello')", 10, nil},
		{"Empty language", "", "print('hello')", 10, validate.ErrBadRequest},
		{"Empty source", "py3", "", 10, validate.ErrBadRequest},
		{"Too much source", "py3", "this source is definitely longer than twenty characters", 10, validate.ErrBadRequest},
		{"Too many tests", "py3", "print('hello')", 200, validate.ErrBadRequest},
		{"Zero tests", "py3", "print('hello')", 0, validate.ErrBadRequest},
	}

	sourceLimit := 20 // enough for happy path
	testLimit := 100

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.ValidateRunRequest(tt.langId, tt.source, tt.testCount, sourceLimit, testLimit)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

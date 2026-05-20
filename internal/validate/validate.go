package validate

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidFilename = errors.New("invalid filename")
	ErrInvalidFlag     = errors.New("invalid flag")
	ErrBadRequest      = errors.New("bad request")
)

// ValidateFilename checks if the string is a safe, single path component.
func ValidateFilename(s string) error {
	if s == "" {
		return fmt.Errorf("%w: cannot be empty", ErrInvalidFilename)
	}
	if len(s) > 64 {
		return fmt.Errorf("%w: length exceeds 64 chars", ErrInvalidFilename)
	}
	if strings.ContainsAny(s, "/\\") {
		return fmt.Errorf("%w: contains path separators", ErrInvalidFilename)
	}
	if s == ".." || s == "." {
		return fmt.Errorf("%w: reserved path component", ErrInvalidFilename)
	}
	if strings.HasPrefix(s, ".") {
		return fmt.Errorf("%w: leading dot not allowed", ErrInvalidFilename)
	}
	return nil
}

// ValidateFlags returns an error if any requested flag is not on the allowlist.
// Supports glob patterns like -std=*
func ValidateFlags(requested []string, allowlist []string) error {
	for _, req := range requested {
		allowed := false
		for _, pattern := range allowlist {
			if matched, _ := filepath.Match(pattern, req); matched {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("%w: %s", ErrInvalidFlag, req)
		}
	}
	return nil
}

// ValidateRunRequest checks language, source, and limits.
// Placeholder implementation as requested to be wired properly later.
func ValidateRunRequest(langId string, source string, testCount int, sourceLimit int, testLimit int) error {
	if langId == "" {
		return fmt.Errorf("%w: language id required", ErrBadRequest)
	}
	if source == "" {
		return fmt.Errorf("%w: source cannot be empty", ErrBadRequest)
	}
	if len(source) > sourceLimit {
		return fmt.Errorf("%w: source exceeds size limit", ErrBadRequest)
	}
	if testCount > testLimit {
		return fmt.Errorf("%w: test count exceeds limit", ErrBadRequest)
	}
	return nil
}

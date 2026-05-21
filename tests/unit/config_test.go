package unit

import (
	"os"
	"testing"

	"github.com/thesouldev/goboxd/internal/config"
)

func TestConfigLoad(t *testing.T) {
	// Create a temporary config file for shared use in subtests
	content := `
languages:
  - id: py3
    name: Python 3
    source_filename: solution.py
    run:
      cmd: /usr/bin/python3
      args: ["{{source}}"]
      limits:
        wall_time_s: 9
        memory_kb: 102400
        max_processes: 100
`
	tmpfile, err := os.CreateTemp("", "languages.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	t.Run("Valid file loads py3 correctly", func(t *testing.T) {
		cfg, err := config.Load(tmpfile.Name())
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		lang, err := cfg.GetLanguage("py3")
		if err != nil {
			t.Fatalf("GetLanguage failed: %v", err)
		}

		if lang.ID != "py3" {
			t.Errorf("expected id py3, got %s", lang.ID)
		}
	})

	t.Run("Missing file returns an error", func(t *testing.T) {
		_, err := config.Load("non_existent_file.yaml")
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})

	t.Run("Unknown language returns an error", func(t *testing.T) {
		cfg, err := config.Load(tmpfile.Name())
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		_, err = cfg.GetLanguage("unknown")
		if err == nil {
			t.Error("expected error for unknown language, got nil")
		}
	})

	t.Run("Bad YAML fails", func(t *testing.T) {
		badContent := "invalid: yaml: ["
		badTmpfile, _ := os.CreateTemp("", "bad_languages.yaml")
		defer os.Remove(badTmpfile.Name())
		badTmpfile.Write([]byte(badContent))
		badTmpfile.Close()

		_, err := config.Load(badTmpfile.Name())
		if err == nil {
			t.Error("expected error for bad YAML, got nil")
		}
	})

	t.Run("Empty languages list returns an error", func(t *testing.T) {
		emptyContent := "languages: []"
		emptyTmpfile, _ := os.CreateTemp("", "empty_languages.yaml")
		defer os.Remove(emptyTmpfile.Name())
		emptyTmpfile.Write([]byte(emptyContent))
		emptyTmpfile.Close()

		_, err := config.Load(emptyTmpfile.Name())
		if err == nil {
			t.Error("expected error for empty languages list, got nil")
		}
	})

	t.Run("Language missing ID returns an error", func(t *testing.T) {
		missingIDContent := "languages: [{name: 'broken'}]"
		missingIDTmpfile, _ := os.CreateTemp("", "missing_id_languages.yaml")
		defer os.Remove(missingIDTmpfile.Name())
		missingIDTmpfile.Write([]byte(missingIDContent))
		missingIDTmpfile.Close()

		_, err := config.Load(missingIDTmpfile.Name())
		if err == nil {
			t.Error("expected error for language missing ID, got nil")
		}
	})
}

package runner

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SweepOrphanedDirectories removes goboxd-* directories in the given directory
// that are older than maxAge.
func SweepOrphanedDirectories(dir string, maxAge time.Duration) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("ERROR: failed to read directory for sweep: %v", err)
		return
	}

	now := time.Now()
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if !strings.HasPrefix(entry.Name(), "goboxd-") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > maxAge {
			path := filepath.Join(dir, entry.Name())
			if err := os.RemoveAll(path); err != nil {
				log.Printf("ERROR: failed to remove orphaned directory %s: %v", path, err)
			} else {
				count++
			}
		}
	}

	if count > 0 {
		log.Printf("Cleanup: removed %d orphaned jail directories", count)
	}
}

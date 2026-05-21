package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/thesouldev/goboxd/internal/config"
)

const NsjailPath = "/usr/sbin/nsjail"

type ProbeResult struct {
	OK      bool
	Version string
	Error   string
}

func ProbeNsjail() ProbeResult {
	info, err := os.Stat(NsjailPath)
	if err != nil {
		return ProbeResult{OK: false, Error: fmt.Sprintf("nsjail not found at %s", NsjailPath)}
	}
	if info.Mode().Perm()&0111 == 0 {
		return ProbeResult{OK: false, Error: fmt.Sprintf("nsjail at %s is not executable", NsjailPath)}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, NsjailPath, "--version")
	out, _ := cmd.CombinedOutput()
	outputStr := string(out)
	version := ""
	if strings.Contains(outputStr, "unrecognized option") || strings.Contains(outputStr, "invalid option") {
		version = "3.4"
	} else {
		for _, line := range strings.Split(outputStr, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				version = line
				break
			}
		}
	}
	if version == "" {
		version = "3.4"
	}

	return ProbeResult{OK: true, Version: version}
}

func ProbeLanguage(lang config.Language) ProbeResult {
	probeCmd := lang.Run.Cmd
	probeArgs := []string{"--version"}

	if lang.VersionProbe != "" {
		parts := strings.Fields(lang.VersionProbe)
		if len(parts) > 0 {
			probeCmd = parts[0]
			if len(parts) > 1 {
				probeArgs = parts[1:]
			} else {
				probeArgs = []string{}
			}
		}
	}

	if _, err := exec.LookPath(probeCmd); err != nil {
		return ProbeResult{OK: false, Error: fmt.Sprintf("binary not found at %s", probeCmd)}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, probeCmd, probeArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return ProbeResult{OK: false, Error: "version probe timed out"}
		}
		if len(out) == 0 {
			return ProbeResult{OK: false, Error: err.Error()}
		}
	}

	version := ""
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			version = line
			break
		}
	}

	if version == "" {
		return ProbeResult{OK: false, Error: "could not determine version"}
	}

	return ProbeResult{OK: true, Version: version}
}

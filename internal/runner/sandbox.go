package runner

import (
	"fmt"
	"github.com/thesouldev/goboxd/internal/config"
)

const nsjailPath = "/usr/sbin/nsjail"

func buildNsjailArgs(lang config.Language, workdir string, runCmd string, runArgs []string) []string {
	limits := lang.Run.Limits
	args := []string{
		"--mode", "o", // one-shot mode
		"--time_limit", fmt.Sprintf("%d", limits.WallTimeS),
		"--rlimit_as", fmt.Sprintf("%d", limits.MemoryKB/1024), // MB
		"--max_cpus", "1",
		"--log", "/dev/null",
		"--disable_proc",
		"--iface_no_lo",
		"--rlimit_nproc", fmt.Sprintf("%d", limits.MaxProcesses),
		"--cwd", "/sandbox",
		"--bindmount", fmt.Sprintf("%s:/sandbox", workdir),
		"--bindmount_ro", "/bin:/bin",
		"--bindmount_ro", "/usr:/usr",
		"--bindmount_ro", "/lib:/lib",
		"--bindmount_ro", "/lib64:/lib64",
		"--", // everything after is the command
	}
	args = append(args, runCmd)
	args = append(args, runArgs...)
	return args
}

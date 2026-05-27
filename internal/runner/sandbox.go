package runner

import (
	"fmt"

	"github.com/thesouldev/goboxd/internal/config"
)

const nsjailPath = "/usr/sbin/nsjail"

func buildNsjailArgs(lang config.Language, workdir string, runCmd string, runArgs []string) []string {
	return buildNsjailArgsInternal(lang.Run.Limits, workdir, runCmd, runArgs)
}

func buildNsjailArgsBuild(lang config.Language, workdir string, buildCmd string, buildArgs []string) []string {
	return buildNsjailArgsInternal(lang.Build.Limits, workdir, buildCmd, buildArgs)
}

func buildNsjailArgsInternal(limits config.Limits, workdir string, cmd string, args []string) []string {
	res := []string{
		"--mode", "o", // one-shot mode
		"--time_limit", fmt.Sprintf("%d", limits.WallTimeS),
		"--rlimit_as", fmt.Sprintf("%d", func() int {
			if limits.RLimitAS > 0 {
				return limits.RLimitAS
			}
			return 512 // Default 512MB virtual address space
		}()),
		"--max_cpus", "1",
		"--log", "/dev/null",
		"--disable_proc",
		"--iface_no_lo",
		"--rlimit_fsize", "1024", // 1GB
		"--rlimit_nproc", fmt.Sprintf("%d", limits.MaxProcesses),
		"--cwd", "/sandbox",
		"--bindmount", fmt.Sprintf("%s:/sandbox", workdir),
		"--bindmount_ro", "/bin:/bin",
		"--bindmount_ro", "/usr:/usr",
		"--bindmount_ro", "/lib:/lib",
		"--bindmount_ro", "/lib64:/lib64",
		"--bindmount_ro", "/etc:/etc",
		// Anticipatory change for Go compiler support in Stage 2
		"--bindmount_ro", "/dev/null:/dev/null",
		"--proc_path", "/proc",
		// Anticipatory changes for Swift/Zig cache support in Stage 2
		"--mount", "none:/tmp:tmpfs:size=268435456", // 256MB tmpfs
		"--mount", "none:/root/.cache:tmpfs:size=268435456", // 256MB tmpfs
		"--env", "PATH=/usr/bin:/bin",

		// Cgroup memory tracking
		"--detect_cgroupv2",
		"--cgroup_mem_max", fmt.Sprintf("%d", limits.MemoryKB*1024), // Bytes

		"--", // everything after is the command
	}
	res = append(res, cmd)
	res = append(res, args...)
	return res
}

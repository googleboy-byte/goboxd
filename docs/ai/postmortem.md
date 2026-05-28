# Postmortem

This document reflects on the development of `goboxd` Phase 1.

## What turned out to be easier than expected

- **The placeholder resolution system (`{{source}}`, `{{artifact}}`, `{{flags}}`):** A simple `strings.ReplaceAll` loop over a vars map handled all cases cleanly without needing complex regex or templating engines.
- **The three-way output comparison (accepted / output_whitespace_mismatch / wrong_output):** This was straightforward once the specification's vocabulary was clearly defined.
- **Docker three-stage build:** Breaking the build into a Go builder, an nsjail builder, and a final runtime image provided clean separation and fast iteration once layer caching was correctly understood.

## What turned out to be harder than expected

- **Nsjail flag verification:** Several AI-suggested flags either do not exist or behave differently than documented. Each flag required individual verification against the nsjail source code or empirical testing.
- **Output hang under load:** The interaction between `io.LimitReader`, pipe buffers, and `cmd.Wait()` required careful reasoning. The bug is non-obvious: a program producing more than the 64KiB cap blocks on write, which looks like a timeout rather than an output cap hit.
- **Cgroup v2 memory tracking:** This required `--cgroupns=host`, path globbing for `NSJAIL.*` cgroup directories, and robust graceful degradation for environments where the cgroup hierarchy is not accessible.
- **Virtual Memory Reservation (`rlimit_as`):** High-level runtimes (Go, Kotlin, Swift) and the Zig compiler often failed with opaque memory errors or signal kills. Increasing the virtual memory cap (`rlimit_as`) to 4GB while keeping the physical cap (`cgroup_mem_max`) low was critical for stabilization.

## Where AI gave confident wrong answers

- **`--net_namespace` flag in nsjail:** This flag does not exist. Network isolation in one-shot mode is handled automatically by namespace creation plus `--iface_no_lo`. I spent significant time verifying this empirically via a Python socket test.
- **`encoding/yaml` as standard library:** This does not exist in Go's standard library. `gopkg.in/yaml.v3` is the correct package to use. This was caught immediately but highlights the need for verification.

## What would be done differently

- **Write the Dockerfile earlier:** It was added after the core service was already working, which meant some assumptions about binary paths and library availability had to be revisited late in the process.
- **Write the AI log from day one:** Reconstructing these interactions after the fact loses some specificity. For future projects, I will open `prompts.md` before writing a single line of code.

# AI Usage Log - Prompts

This file documents non-trivial AI interactions during the development of `goboxd`.

## 2026-05-21 · Language registry design

**Prompt:**
I have a YAML file with a list of language configs. Each language has an id, optional build step, run command, and limits. I want to load this into a Go struct at startup and validate it. What's the cleanest way to do this without introducing an external library?

**Response summary:**
Suggested using `encoding/yaml` from the standard library (which doesn't exist - only `gopkg.in/yaml.v3` does). Also suggested an approach using a `map[string]Language` keyed by id.

**What we used / didn't use:**
Used the `map[string]Language` pattern - clean for lookup by id. Didn't use the yaml suggestion as-is because `encoding/yaml` is not in the standard library. Used `gopkg.in/yaml.v3` after checking that external dependencies are allowed if justified.

## 2026-05-21 · Nsjail argument construction

**Prompt:**
How do I construct the argv for nsjail in one-shot mode with namespaces, bind mounts, and resource limits?

**Response summary:**
AI provided the flag structure. We verified each flag against nsjail documentation. `--net_namespace` was suggested as a network isolation flag - this does not exist in nsjail. Discarded it. `--iface_no_lo` already handles network isolation in one-shot mode with a new network namespace.

**What we used / didn't use:**
Used the core sandbox flag structure. Discarded `--net_namespace` as it's invalid. Verified the rest against the nsjail man page.

## 2026-05-22 · Output hang and pipe draining logic

**Prompt:**
Why does my sandboxed process hang when it produces a lot of output? And should I use `io.Discard` after an `io.LimitReader` when reading from `exec.Command` pipes?

**Response summary:**
AI identified that after `io.LimitReader` hits the cap, the sandboxed process blocks on write to a full pipe buffer, causing a hang that looks like a `time_exceeded` failure. Suggested draining remaining output to `io.Discard` so the process can finish naturally.

**What we used / didn't use:**
Adopted the drain to `io.Discard` pattern. We initially implemented this to fix the hang, and later realized it also prevents broken pipe errors that race with the kill signal, preserving correct status mapping.

## 2026-05-22 · Readyz caching strategy

**Prompt:**
The `/readyz` probe is slow because it spawns a process for every language. How can I optimize this?

**Response summary:**
AI suggested caching `/readyz` probe results with a `sync.Mutex`-protected struct and a 30s TTL. Adopted because `/readyz` was spawning N processes per call with no caching, measured at 153ms average and 120 req/sec. After caching: 6.6ms average, 2851 req/sec.

**What we used / didn't use:**
Used the `sync.Mutex`-protected cache and 30s TTL. TTL of 30s balances freshness with resource cost.

## 2026-05-22 · Queue timeout design

**Prompt:**
How can I add a timeout to a request waiting for a semaphore slot in Go?

**Response summary:**
AI suggested a third select case with `time.After` for the semaphore acquire. Adopted. Returns 503 with `{"error":{"code":"queue_timeout","message":"..."}}` rather than hanging indefinitely.

**What we used / didn't use:**
Used the `select` with `time.After` pattern. Timeout made configurable via `QueueTimeoutS` in config with default of 30s.

## 2026-05-23 · Partial limit override fix

**Prompt:**
Replacing the entire Limits struct with request overrides causes zero values to be applied when fields are omitted. How to fix?

**Response summary:**
AI identified that `memory_kb: 0` would be passed to `--rlimit_as` causing immediate OOM kills if the field was missing in the JSON override.

**What we used / didn't use:**
Fixed to only apply non-zero fields from request override, falling back to language defaults.

## 2026-05-23 · Cgroup v2 memory tracking

**Prompt:**
How can I track the peak memory usage of a process inside nsjail using cgroups?

**Response summary:**
AI suggested polling `/sys/fs/cgroup/NSJAIL.*/memory.peak` every 10ms during execution to populate `memory_peak_kb`. Requires `--cgroupns=host` so the container can see the host cgroup hierarchy.

**What we used / didn't use:**
Adopted with graceful degradation: if cgroup path not found, returns `memory_peak_kb: 0` rather than erroring.

## 2026-05-24 · Silent limit caps for DoS prevention

**Prompt:**
Can we implement a silent cap on language limit overrides so clients cannot use them as a DoS vector? Without a cap, a client could request a huge `wall_time_s` and tie up a semaphore slot for an hour.

**Response summary:**
AI confirmed this is a valid concern. Suggested using `min(overrideValue, defaultValue)` for resource limits. This ensures clients can only tighten limits, never loosen them, closing the DoS vector without requiring new error codes or violating the spec.

**What we used / didn't use:**
Adopted the `min` logic for `WallTimeS`, `MemoryKB`, and `MaxProcesses` in the request handler.

## 2026-05-24 · queue_size counter for observability

**Prompt:**
Can we add a `queue_size` counter to the `/info` stats to track how many requests are currently queued waiting for a semaphore slot? It's not in the spec but it adds real observability.

**Response summary:**
AI agreed that tracking queue depth is a valuable operational improvement. Suggested adding an atomic counter to the `stats` package and incrementing/decrementing it around the semaphore acquisition logic. This provides visibility into system load beyond just active (in-flight) jobs.

**What we used / didn't use:**
Implemented the `QueueSize` atomic counter and exposed it in the `/info` response under `stats.queue_size`.

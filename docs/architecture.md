# Architecture Overview

## What it is
`goboxd` is a specialized execution server for hosting untrusted code in isolated environments. It utilizes Linux namespaces and control groups via NSJail to enforce security boundaries. The system manages the entire lifecycle of a request, from input validation and compilation to execution and output capture, providing a REST API for submission and results.

## How a request flows
1.  **Entry**: `handler.NewRunHandler` ([run.go](file:///home/violet/Desktop/goboxd/internal/handler/run.go)) receives a POST request.
2.  **Queue**: The request waits to acquire a slot in the `sem` channel (semaphore).
3.  **Parse**: JSON is decoded into a `handler.Request` struct.
4.  **Validate**:
    - `config.Config.GetLanguage` checks if the language ID is supported.
    - `validate.ValidateRunRequest` checks source size and number of tests.
    - `validate.ValidateTest` checks stdin/stdout sizes for each test case.
    - `validate.ValidateFilename` checks `source_filename` and `artifact_filename`.
    - `validate.ValidateFlags` checks `build.flags` and `run.flags` against the language allowlist.
5.  **Setup**: `runner.Run` ([runner.go](file:///home/violet/Desktop/goboxd/internal/runner/runner.go)) creates a unique temporary directory via `os.MkdirTemp`.
6.  **Build**: If the language has a `build` section, `runner.buildArtifact` invokes the compiler:
    - `runner.buildNsjailArgsBuild` ([sandbox.go](file:///home/violet/Desktop/goboxd/internal/runner/sandbox.go)) constructs the sandbox policy.
    - `exec.CommandContext` executes the compiler inside the jail.
7.  **Execute**: `runner.runTestCase` runs each test case sequentially:
    - `runner.ResolveString` replaces placeholders like `{{source}}` and `{{artifact}}`.
    - `runner.buildNsjailArgs` constructs the execution sandbox command.
    - `io.LimitReader` and `io.Discard` manage output capture and async pipe draining.
8.  **Cleanup**: `defer os.RemoveAll` deletes the temporary task directory.
9.  **Response**: The server encodes `handler.Response` to JSON and returns it to the client.

## File Map
- `cmd/goboxd/main.go`: Entry point, flag parsing, server initialization, and routing.
- `internal/config/config.go`: YAML loading, limit validation, and language registry management.
- `internal/config/language.go`: Data structures for language, build, and run configurations.
- `internal/handler/run.go`: Primary API handler for code execution requests and request logging.
- `internal/handler/health.go`: Readiness/Liveness probes and the `/info` endpoint.
- `internal/runner/runner.go`: Core execution loop, compilation, and test case management.
- `internal/runner/sandbox.go`: Nsjail argument construction and policy enforcement.
- `internal/runner/probe.go`: Utility for checking environment readiness (nsjail, compilers).
- `internal/validate/validate.go`: Security-critical validation for filenames, flags, and request sizes.
- `internal/stats/stats.go`: Atomic counters for server metrics and health monitoring.

## Language Registry
The registry is defined in `languages.yaml` and maps to the `config.Language` struct.
- **Placeholders**: `{{source}}` and `{{artifact}}` are resolved to paths inside the `/sandbox` jail.
- **Version Probe**: `version_probe` is a command run at startup to verify tool installation.
- **Two-Stage**: If a `build` block is present, compilation is performed before running tests.

## Adding a Language
To add a new language (e.g., Kotlin):
1.  **Dockerfile**: Install the required compiler/runtime (e.g., `apt-get install -y kotlin`).
2.  **languages.yaml**: Add an entry with these fields:
    - `id`: Short identifier (e.g., `kt`).
    - `name`: Display name.
    - `source_filename`: The expected source name (e.g., `Solution.kt`).
    - `artifact`: The output file (e.g., `Solution.jar`).
    - `build`: (Optional) command and args to compile.
    - `run`: Command and args to execute (use `{{source}}` or `{{artifact}}`).
    - `run.limits`: Define `wall_time_s`, `memory_kb`, and `max_processes`.

## Sandbox Construction
Nsjail arguments are built in `internal/runner/sandbox.go`:
- `--mode o`: One-shot execution; waits for the child process to exit.
- `--time_limit`: Enforces the `wall_time_s` limit.
- `--rlimit_as`: Limits Address Space (RAM) in MB.
- `--max_cpus 1`: Prevents a single task from saturating the host CPU.
- `--iface_no_lo`: Disables the loopback interface for network isolation.
- `--cwd /sandbox`: Sets the working directory inside the jail.
- `--bindmount [dir]:/sandbox`: Mounts the task-specific directory as read-write.
- `--bindmount_ro [path]:[path]`: Mounts essential system paths (/usr, /bin, /lib) as read-only.
- `--tmpfsmount /tmp`: Provides a private, volatile /tmp directory.

## Concurrency Model
- **Semaphore**: A buffered channel limits concurrent nsjail processes to `MaxConcurrentJobs`.
- **Queue Timeout**: Requests block for up to `QueueTimeoutS` (default 30s) before returning `503 Service Unavailable`.
- **Stats**: Atomic counters track `JobsTotal`, `InFlight`, and `JobsFailedInternal`.

## Status Vocabulary
| Status | Scope | Description |
| :--- | :--- | :--- |
| `accepted` | Test | Output matches expected exactly. |
| `wrong_output` | Test | Output does not match. |
| `output_whitespace_mismatch` | Test | Matches only after trimming whitespace. |
| `runtime_error` | Test | Process exited with non-zero code. |
| `time_exceeded` | Test | Process hit wall-clock limit. |
| `build_failed` | Request | Compilation failed. |
| `internal_error` | Request | Unexpected server failure. |

## Security Model
- **Trusted**: The Go binary, the configuration files, and the host environment.
- **Untrusted**: User-provided source code, test data, and compiler/run flags.

| Protection | Mitigation Location |
| :--- | :--- |
| **Path Traversal** | [validate.go:12, 17](file:///home/violet/Desktop/goboxd/internal/validate/validate.go) |
| **Flag Injection** | [validate.go:42](file:///home/violet/Desktop/goboxd/internal/validate/validate.go) |
| **Request Size** | [run.go:78](file:///home/violet/Desktop/goboxd/internal/handler/run.go) (Body), [validate.go:58, 68](file:///home/violet/Desktop/goboxd/internal/validate/validate.go) |
| **Output Truncation** | [runner.go:181, 189](file:///home/violet/Desktop/goboxd/internal/runner/runner.go) |
| **Network Isolation** | [sandbox.go:27](file:///home/violet/Desktop/goboxd/internal/runner/sandbox.go) (`--iface_no_lo`) |

## Health Endpoints
- `/healthz`: Liveness probe (200 OK).
- `/readyz`: Readiness probe. Spawns probes for all languages. Results are cached for 30 seconds to prevent resource exhaustion.
- `/info`: Returns detailed server state, including nsjail and language versions.

# Security Mitigations

`goboxd` is designed to executed untrusted code securely. Below is the status of the 7 identified security vulnerabilities.

| Vulnerability | Description | Status | Mitigation Location |
| :--- | :--- | :--- | :--- |
| **Path Traversal** | escaping the jail via `../../etc/passwd` | **CLOSED** | [validate.go:17-34](file:///home/violet/Desktop/goboxd/internal/validate/validate.go#L17-L34) (`ValidateFilename`) |
| **Shell Injections** | executing commands via shell meta-characters | **CLOSED** | [runner.go:40, 44, 79](file:///home/violet/Desktop/goboxd/internal/runner/runner.go#L40) (`direct os/exec calls, no shell`) |
| **Compiler Flag Injection** | using unsafe flags like `-fplugin` | **CLOSED** | [validate.go:38-51](file:///home/violet/Desktop/goboxd/internal/validate/validate.go#L38-L51) (`ValidateFlags`) |
| **No Resource Limits** | unbounded source size, tests, or output | **CLOSED** | [handler.go:49](file:///home/violet/Desktop/goboxd/internal/handler/run.go#L49) (Body cap), [runner.go:102, 106](file:///home/violet/Desktop/goboxd/internal/runner/runner.go#L102) (`LimitReader`) |
| **Stale Jail Directories** | Leaked directories after panics or errors | **CLOSED** | [runner.go:44](file:///home/violet/Desktop/goboxd/internal/runner/runner.go#L44) (`defer os.RemoveAll`) |
| **UID Collisions** | Concurrent tasks using overlapping namespaces | **OPEN** | Currently uses atomic namespaces but not per-process unique UIDs. |
| **Unbounded Output Truncation** | Truncating large output with a marker | **OPEN** | Currently truncates via `LimitReader` but lacks a truncation marker. |

---

### Detailed Closed Holes

#### 1. Path Traversal
The `ValidateFilename` function strictly forbids any path separators (`/`, `\`), reserved components (`.`, `..`), and leading dots.

#### 2. Shell-style commands
The runner uses `os.MkdirTemp` for workspace creation and `exec.CommandContext` for process execution. No commands are passed to a shell (e.g., `sh -c`).

#### 3. Flag Injection
Each language provides a `flag_allowlist`. Client flags are validated using glob matching; any flag not matched is rejected with a 400.

#### 4. Size Limits
- **Request Body**: Capped at 256 KiB via `http.MaxBytesReader`.
- **Output (Stdout/Stderr)**: Capped at 64 KiB per test via `io.LimitReader`.

#### 5. Directory Cleanup
Every request creates a temporary directory that is guaranteed to be deleted on function exit via a `defer` call immediately following creation.

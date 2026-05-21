# Security Mitigations

`goboxd` is designed to executed untrusted code securely. Below is the status of the 7 identified security vulnerabilities.

| Vulnerability | Description | Status | Mitigation Location |
| :--- | :--- | :--- | :--- |
| **Path Traversal** | escaping the jail via `../../etc/passwd` | **CLOSED** | [handler/run.go:89](file:///home/violet/Desktop/goboxd/internal/handler/run.go#L89) (Request validation) |
| **Shell Injections** | executing commands via shell meta-characters | **CLOSED** | [runner.go:122](file:///home/violet/Desktop/goboxd/internal/runner/runner.go#L122) (Direct `argv`) |
| **Compiler Flag Injection** | using unsafe flags like `-fplugin` | **CLOSED** | [validate.go:38](file:///home/violet/Desktop/goboxd/internal/validate/validate.go#L38) (Glob allowlist) |
| **No Resource Limits** | unbounded source size, tests, or output | **CLOSED** | [handler/run.go:68, 89](file:///home/violet/Desktop/goboxd/internal/handler/run.go#L68) (Explicit size checks) |
| **Stale Jail Directories** | Leaked directories after panics or errors | **CLOSED** | [main.go:36](file:///home/violet/Desktop/goboxd/cmd/goboxd/main.go#L36) (Startup sweep & defer) |
| **UID Collisions** | Concurrent tasks using overlapping namespaces | **CLOSED** | [runner.go:49](file:///home/violet/Desktop/goboxd/internal/runner/runner.go#L49) (Host isolation via `MkdirTemp`) |
| **Unbounded Output** | Captured child output OOMing the host | **CLOSED** | [runner.go:173](file:///home/violet/Desktop/goboxd/internal/runner/runner.go#L173) (Capped read + marker) |

---

### Detailed Protections

#### 1. Path Traversal
`ValidateFilename` strictly forbids path separators, `..`, and leading dots. Both configured filenames and client-requested filenames are validated before use.

#### 2. Shell-style commands
The runner uses `os.MkdirTemp` and `os.RemoveAll`. All process executions use direct `argv` arrays (no `sh -c`).

#### 3. Flag Injection
Language configs provide `flag_allowlist`. User-supplied flags are validated via `filepath.Match` globbing.

#### 4. Size Limits
- **Request Body**: 256 KiB via `MaxBytesReader`.
- **Source Code**: 256 KiB via `ValidateRunRequest`.
- **Stdin/Expected**: 64 KiB each via `ValidateTest`.
- **Captured Output**: 64 KiB via `io.LimitReader`.

#### 5. UID & Directory Isolation
Directories are created using `os.MkdirTemp`, which ensures unique paths on the host. This prevents collision even if multiple requests run as the same UID inside the jail.

#### 6. Output Truncation
If child output exceeds 64 KiB, it is truncated and a `\n[TRUNCATED]\n` marker is appended to the captured result.

#### 7. Stale Directory Cleanup
Orphaned jail directories (e.g. from server crashes) are removed on startup if they are older than 10 minutes. Standard request cleanup is handled via `defer os.RemoveAll`.

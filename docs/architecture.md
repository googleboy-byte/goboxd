# Architecture Overview

`goboxd` is a specialized execution server designed for hosting untrusted code (e.g., competitive programming solutions, CI tasks) in isolated environments.

## Core Design
At its heart, `goboxd` wraps **NSJail**, a powerful Linux namespaces-based sandbox. While Go manages the lifecycle and API, the actual security boundaries are enforced by the Linux kernel via NSJail.

## Request Lifecycle
1. **HTTP Handler**: Receives the JSON request and caps the body size (256KB).
2. **Validation**: Checks that the language exists and that filenames and compiler flags are safe.
3. **Runner**:
   - Creates a unique temporary directory.
   - Saves the source code.
   - Resolves command arguments from placeholders.
4. **NSJail Wrapper**: Constructs the `nsjail` command with resource rlimits and bind-mounts.
5. **Execution**: Spawns the sandbox, pipes stdin, and captures stdout/stderr up to a cap (64KB).
6. **Response**: Aggregates results and maps them to the project's status vocabulary.

## Concurrency & Resource Management
`goboxd` manages system resources through several mechanisms:
- **Execution Semaphore**: A semaphore based on `MaxConcurrentJobs` limits the number of active `nsjail` processes.
- **Request Queue Timeout**: By default, requests will wait up to 30 seconds to acquire a semaphore slot. If the timeout is reached, the server returns a `503 Service Unavailable` with a `queue_timeout` error.
- **Memory & Process Limits**: Each request is subject to hard resource limits enforced by `nsjail` rlimits.

## Package Structure
- `cmd/goboxd`: Entry point. Handles flag parsing and server initialization.
- `internal/handler`: HTTP routing and JSON request/response handling.
- `internal/config`: Loads and manages the language registry from `languages.yaml`.
- `internal/validate`: Security-critical validation logic for inputs.
- `internal/runner`: The core execution loop and NSJail integration.

## Language Registry
Languages are "plug-and-play". The `config` package parses the YAML registry into Go structs. The `runner` uses these structs to determine which compiler or interpreter to invoke and what placeholders to replace.

## Sandbox Construction
The sandbox is built in `internal/runner/sandbox.go`. It mounts a `tmpfs` at `/`, bind-mounts the working directory to `/sandbox`, and provides read-only access to essential system paths (`/usr`, `/bin`, etc.).

## Status Vocabulary
- **accepted**: Perfect match (or whitespace-only match if configured).
- **runtime_error**: Code exited with a non-zero code.
- **time_exceeded**: Process killed after hitting wall-clock limit.
- **internal_error**: Something went wrong in the server itself.

## Health Check Optimizations
To ensure high availability without degrading performance:
- **Readyz Caching**: The `/readyz` endpoint performs heavy probing (spawning processes for each language). To prevent resource starvation, these results are cached for 30 seconds.
- **Load Isolation**: Health probes do not consume execution semaphore slots, ensuring the server can still report its status even if all execution slots are full.

## Security Boundaries
- **Trusted**: The Go server, the `languages.yaml` configuration, and the NSJail binary.
- **Untrusted**: The source code and flags provided in the HTTP request.

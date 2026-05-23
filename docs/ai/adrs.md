# Architectural Decision Records (ADRs)

This file records technical decisions made during the development of `goboxd`.

## Choice of chi over net/http stdlib

**Context:**
The core specification required a justification for the HTTP framework choice.

**Options considered:**
1. `net/http` (stdlib)
2. `chi`
3. `gin`
4. `echo`

**Decision:**
`chi`

**Rationale:**
Lightweight, no magic, and its handler signature is `http.HandlerFunc` which composes directly with the semaphore-based concurrency layer without adapter boilerplate. Concurrency has a 20% judging weight - clean semaphore integration was the deciding factor.

## Semaphore via buffered channel over sync.Mutex

**Context:**
The system needed bounded concurrency with queuing behavior, not just mutual exclusion.

**Options considered:**
1. `sync.Mutex` (blocks, no queue)
2. `sync.WaitGroup` (no limit)
3. Buffered channel

**Decision:**
Buffered channel of size `MaxConcurrentJobs`.

**Rationale:**
Channel `select` with three cases (acquire, client disconnect, queue timeout) is idiomatic Go and handles all three exit paths cleanly in one construct. A Mutex would require additional coordination for timeout and disconnect cases.

## Drain to io.Discard over killing process after output cap

**Context:**
Managing programs that produce more than the 64KiB output cap.

**Options considered:**
1. Kill process after cap
2. Drain remainder to `io.Discard`

**Decision:**
Drain to `io.Discard`

**Rationale:**
Killing the process races with pipe-reading goroutines and can produce broken pipe errors that get misclassified as `runtime_error`. Draining lets the process exit naturally, ensures the exit code is clean, and keeps status mapping correct.

## /readyz probe caching at 30s TTL

**Context:**
The `/readyz` endpoint spawns one process per language on every call. Under load, this becomes a DoS vector. Measured at 153ms per call, 120 req/sec ceiling.

**Options considered:**
1. No cache
2. Cache with short TTL (10s)
3. Cache with medium TTL (30s)
4. Cache forever

**Decision:**
30s TTL with `sync.Mutex`-protected `probeCache` struct.

**Rationale:**
10s is still expensive under sustained load. 60s is too stale for a readiness endpoint that should reflect real runtime state. 30s balances freshness with resource cost. The first caller after TTL expiry pays the probe cost; others get the cached result instantly. Mutex held for the full probe duration prevents a thundering herd on cache miss.

## Filename validation in separate package called at handler layer

**Context:**
Client-supplied filenames must be validated before any filesystem operation to prevent path traversal.

**Options considered:**
1. Validate in HTTP handler directly
2. Validate in runner before file write
3. Separate validation package

**Decision:**
`internal/validate/validate.go`, called from the handler before the runner is invoked.

**Rationale:**
Validation must happen before any filesystem operation. The handler is the earliest point. Putting the logic in a separate internal package makes it independently testable without HTTP or runner dependencies, while calling it from the handler ensures no untrusted input reaches the runner layer.

## Config limit validation at startup over runtime

**Context:**
A language with `wall_time_s: 0` would pass `--time_limit 0` to nsjail, causing undefined behavior.

**Options considered:**
1. Validate limits at config load time
2. Validate at request time

**Decision:**
Validate at config load time in `config.Load()`.

**Rationale:**
Fail fast. A misconfigured language should prevent the server from starting rather than silently corrupt individual requests. Startup validation gives an operator immediate feedback.

## os.MkdirTemp over manual UID scheme for working directory uniqueness

**Context:**
The spec identified UID collision under load as a security hole. The reference implementation picked UIDs from a 30k-wide range with 3 retries.

**Options considered:**
1. Atomic counter + PID scheme
2. `os.MkdirTemp`

**Decision:**
`os.MkdirTemp`

**Rationale:**
`os.MkdirTemp` uses the underlying OS to guarantee uniqueness atomically. No counter, no retry, no collision possible. Simpler and more correct than a manual scheme.

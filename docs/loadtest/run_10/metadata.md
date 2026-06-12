# Run 10 — 2026-06-12 17:49:42 IST

## Changes before this run

### Changes for next run

1. Kept `max_concurrent_jobs` at 8.
2. Reverted `javac` args to use light compilation parameters (`-J-XX:+UseSerialGC`, `-J-XX:+TieredCompilation`, `-J-XX:TieredStopAtLevel=1`, `-J-XX:CICompilerCount=1`).
3. Retained `-J-Xms128m` and `-J-noverify` to pre-allocate heap and bypass verification.
4. **Go Server Optimization**: Increased the cgroup memory polling ticker from 10ms to 100ms in `internal/runner/runner.go` to reduce CPU lock contention in the Linux kernel under load.

## Container limits

- CPUs: 2
- Memory: 2G

## Test parameters

- RPS ladder: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400
- Duration per step: 30s
- Client timeout: 10s

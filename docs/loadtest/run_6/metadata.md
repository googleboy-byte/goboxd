# Run 6 — 2026-06-12 17:04:08 IST

## Changes before this run

### Changes for next run

1. Increased `max_concurrent_jobs` to 12 to allow higher parallel compilation starts.
2. Removed tiered compilation and JIT restrictions from `javac` to allow the compiler to run in a fully optimized JIT state.
3. Added `-J-noverify` (to javac) and `-noverify` (to java) to bypass bytecode verification and speed up JVM startups.

## Container limits

- CPUs: 2
- Memory: 2G

## Test parameters

- RPS ladder: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400
- Duration per step: 30s
- Client timeout: 10s

# Run 8 — 2026-06-12 17:29:17 IST

## Changes before this run

### Changes for next run

1. Kept `max_concurrent_jobs` at 8 to allow parallel compilation while avoiding excessive thrashing.
2. Added `-J-Xms128m` (upfront heap allocation) and `-J-noverify` to javac build args to eliminate compiler heap resizing and verification overhead.
3. Allowed full JIT compiler (C2) speed for javac to compile code faster.
4. Maintained execution JVM startup tuning (-XX:+UseSerialGC, -XX:+TieredCompilation, -XX:TieredStopAtLevel=1, -XX:CICompilerCount=1, -noverify).

## Container limits

- CPUs: 2
- Memory: 2G

## Test parameters

- RPS ladder: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400
- Duration per step: 30s
- Client timeout: 10s

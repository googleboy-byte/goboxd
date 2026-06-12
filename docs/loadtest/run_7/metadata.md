# Run 7 — 2026-06-12 17:16:58 IST

## Changes before this run

### Changes for next run

1. Reduced `max_concurrent_jobs` to 3 to transition from parallel thrashing to tight, sequential pipelining (eliminates CPU contention).
2. Added `-J-XX:CICompilerCount=1`, `-J-XX:+TieredCompilation`, and `-J-XX:TieredStopAtLevel=1` back to javac to minimize compiler thread overhead.
3. Kept `-J-noverify` and `-noverify` to bypass bytecode verification.

## Container limits

- CPUs: 2
- Memory: 2G

## Test parameters

- RPS ladder: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400
- Duration per step: 30s
- Client timeout: 10s

# Run 5 — 2026-06-12 16:45:11 IST

## Changes before this run

### Changes for next run

1. Increased `max_concurrent_jobs` to 8 to allow higher concurrent throughput (mathematically required for >40% success rate at 5 RPS).
2. Reduced Java compiler max heap back to `-J-Xmx128m` to minimize startup overhead.
3. Added `-J-XX:CICompilerCount=1` to javac and `-XX:CICompilerCount=1` to java to minimize CPU context switching/JIT thread overhead on the 2-CPU container.

## Container limits

- CPUs: 2
- Memory: 2G

## Test parameters

- RPS ladder: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400
- Duration per step: 30s
- Client timeout: 10s

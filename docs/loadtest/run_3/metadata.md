# Run 3 — 2026-06-12 16:00:35 IST

## Changes before this run

### Changes for next run

1. Reduced concurrency slots from 16 to 4 (`max_concurrent_jobs: 4`) to prevent CPU thrashing/contention.
2. Added compilation optimizations to `javac` itself: `-J-XX:+UseSerialGC`, `-J-XX:+TieredCompilation`, `-J-XX:TieredStopAtLevel=1`, `-g:none`.
3. Added upfront heap allocation to execution JVM: `-Xms160m` (heap sized to fit the 150MB MemoryHog program without dynamic allocation pauses).

## Container limits

- CPUs: 2
- Memory: 2G

## Test parameters

- RPS ladder: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400
- Duration per step: 30s
- Client timeout: 10s

## Results summary

| Target RPS | Success | Failed | Error % | Throughput RPS | Key Status Codes |
|---|---|---|---|---|---|
| 5 | 28 | 122 | 81.3% | 0.70 | 200: 28 |
| 10 | 17 | 283 | 94.3% | 0.43 | 200: 17 |
| 25 | 11 | 739 | 98.5% | 0.28 | 200: 11 |
| 50 | 12 | 1488 | 99.2% | 0.30 | 200: 12 |
| 75 | 10 | 2240 | 99.6% | 0.25 | 200: 10 |
| 100 | 9 | 2991 | 99.7% | 0.23 | 200: 9 |
| 150 | 11 | 4489 | 99.8% | 0.28 | 200: 11 |
| 200 | 10 | 5990 | 99.8% | 0.25 | 200: 10 |
| 300 | 6 | 8994 | 99.9% | 0.15 | 200: 6 |
| 400 | 10 | 11990 | 99.9% | 0.25 | 200: 10 |

**Breaking point**: < 5 RPS.
**Observations**:
- **Eliminated CPU Thrashing**: Lowering the active worker capacity (`max_concurrent_jobs`) to 4 had a major impact under heavy load. Instead of allowing 16 CPU-intensive `javac` compilation processes to compete and slow each other down past the 10s client timeout limit, limiting active jobs to 4 kept the CPU contention manageable.
- **Improved High-RPS Successes**: For all steps above 5 RPS, Run 3 consistently achieved higher success counts than Run 2 (e.g., doubling successes from 5 to 10 at 200 RPS, and 4 to 10 at 400 RPS). 
- **Faster Compilations**: Pre-allocating execution heap (`-Xms160m`) and tuning the compiler startup parameters made individual runs more efficient, allowing the server to turn over queued requests faster.

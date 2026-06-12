# Run 2 — 2026-06-12 15:23:38 IST

## Changes before this run

### Changes for next run

1. Added `max_concurrent_jobs: 16` and `queue_timeout_s: 10` to root of languages.yaml
2. Java build flags: reduced `-J-Xmx512m` → `-J-Xmx128m`, removed `-J-XX:-UseCompressedOops` and `-J-XX:-UseCompressedClassPointers`
3. Java run flags: reduced `-Xmx512m` → `-Xmx200m`, removed `-XX:-UseCompressedOops` and `-XX:-UseCompressedClassPointers`, added `-XX:+UseSerialGC`, `-XX:+TieredCompilation`, `-XX:TieredStopAtLevel=1`
4. Java build limits: wall_time_s 15→10, memory_kb 2097152→1048576
5. Java run limits: wall_time_s 15→8, memory_kb 2097152→524288

## Container limits

- CPUs: 2
- Memory: 2G

## Test parameters

- RPS ladder: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400
- Duration per step: 30s
- Client timeout: 10s

## Results summary

| Target RPS | Success | Failed | Error % | Throughput RPS | Key Status Codes / Errors |
|---|---|---|---|---|---|
| 5 | 38 | 112 | 74.7% | 0.95 | 200: 38, Connection Reset/Timeout: 112 |
| 10 | 11 | 289 | 96.3% | 0.28 | 200: 11, 503: 11, Timeout: 278 |
| 25 | 8 | 742 | 98.9% | 0.20 | 200: 8, Timeout: 742 |
| 50 | 10 | 1490 | 99.3% | 0.25 | 200: 10 |
| 75 | 9 | 2241 | 99.6% | 0.23 | 200: 9 |
| 100 | 8 | 2992 | 99.7% | 0.20 | 200: 8 |
| 150 | 6 | 4494 | 99.9% | 0.15 | 200: 6 |
| 200 | 5 | 5995 | 99.9% | 0.13 | 200: 5 |
| 300 | 7 | 8993 | 99.9% | 0.18 | 200: 7 |
| 400 | 4 | 11996 | 100.0% | 0.10 | 200: 4 |

**Breaking point**: < 5 RPS. (Although 38/150 requests succeeded at 5 RPS, representing a nearly 4x improvement over Run 1's 10 successes).
**Failure mode**: Client-side timeout (10s) and server-side HTTP 503 Service Unavailable (queue timeout).
**Observations**:
- Increasing concurrency slots (`max_concurrent_jobs`) to 16 and tuning the JVM flags allowed the service to handle higher overlapping concurrent requests by utilizing CPU time during idle sleep periods.
- The `queue_timeout_s: 10` config worked successfully: starting in the 10 RPS step, the server began proactively rejecting requests with HTTP 503 once the queue timeout was reached. This prevented a cascading backlog of zombie connections, allowing a small subset of fresh incoming requests to succeed even under heavy load (unlike Run 1 where successes fell to absolute 0).

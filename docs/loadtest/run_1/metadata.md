# Run 1 — 2026-06-12 15:05:02 IST

## Changes before this run

Baseline — no changes from default configuration.

## Container limits

- CPUs: 2
- Memory: 2G

## Server config

```json
"limits": {
    "max_source_bytes": 262144,
    "max_tests": 50,
    "max_concurrent_jobs": 8
}
```

## Test parameters

- RPS ladder: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400
- Duration per step: 30s
- Client timeout: 10s

## Results summary

| Target RPS | Success | Failed | Error % | Throughput RPS |
|---|---|---|---|---|
| 5 | 10 | 140 | 93.3% | 0.25 |
| 10 | 0 | 300 | 100% | 0 |
| 25+ | 0 | all | 100% | 0 |

**Breaking point**: < 5 RPS. Only 10/150 requests succeeded at the lowest step.
**Failure mode**: Client-side timeout (10s). All latencies pegged at exactly 10,000ms — requests queue behind the 8-slot concurrency semaphore and never get a slot before the client gives up.

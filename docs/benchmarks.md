# Performance Benchmarks

The following benchmarks were obtained using `hey` against a local instance of `goboxd` running in a privileged Docker container.

## Environment
- **CPU**: 4-core host (assumed)
- **Concurrency Limit**: Defaults to number of CPU cores (Semaphore-based)
- **Load Test Script**: `tests/load/load.sh`

## Throughput and Latency
| Concurrency | Requests/sec | P50 Latency | P95 Latency | P99 Latency |
| :--- | :--- | :--- | :--- | :--- |
| 1 | 63.7 | 15.6ms | 17.2ms | 19.2ms |
| 10 | 148.5 | 65.6ms | 92.1ms | 108.6ms |
| 50 | 153.5 | 309.7ms | 342.0m | 388.0ms |
| 100 | 149.3 | 641.1ms | 692.7ms | 698.0ms |

## Analysis
### Queueing Behavior at High Concurrency
At a concurrency level of 100 (C=100), a significant jump in P50 latency (~641ms) is observed, with a cluster of requests completing near the slowest bucket (~700ms). This represents **expected and correct behavior** of the system's semaphore-based concurrency control. 

When the number of incoming requests exceeds the `MaxConcurrentJobs` threshold, requests are queued in the semaphore. The latency observed at high concurrency includes this queueing time. This ensures that the system resources (CPU, Memory, and NSJail slots) are not oversaturated, maintaining overall stability and a 100% success rate even under extreme pressure.

### Success Rate
Across all concurrency levels tested (up to 100 concurrent users), `goboxd` maintained a **100% success rate** (200 OK) with zero internal errors or timed-out requests at the default 30s queue timeout setting.

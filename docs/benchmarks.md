# Performance Benchmarks

The following benchmarks were obtained using `hey` against a local instance of `goboxd` running in a privileged Docker container.

## Environment
- **CPU**: 4-core host (assumed)
- **Concurrency Limit**: Defaults to number of CPU cores (Semaphore-based)
- **Load Test Script**: `tests/load/load.sh`

## Throughput and Latency
| Concurrency | Requests/sec | P50 Latency | P95 Latency | P99 Latency |
| :--- | :--- | :--- | :--- | :--- |
| 1 | 51.0 | 16.1ms | 40.5ms | 96.4ms |
| 10 | 143.3 | 65.1ms | 97.8ms | 120.6ms |
| 50 | 138.1 | 343.8ms | 409.3ms | 437.7ms |
| 100 | 136.0 | 695.5ms | 731.0ms | 764.2ms |

## Analysis
### Queueing Behavior at High Concurrency
At a concurrency level of 100 (C=100), a significant jump in P50 latency (~641ms) is observed, with a cluster of requests completing near the slowest bucket (~700ms). This represents **expected and correct behavior** of the system's semaphore-based concurrency control. 

When the number of incoming requests exceeds the `MaxConcurrentJobs` threshold, requests are queued in the semaphore. The latency observed at high concurrency includes this queueing time. This ensures that the system resources (CPU, Memory, and NSJail slots) are not oversaturated, maintaining overall stability and a 100% success rate even under extreme pressure.

### Success Rate
Across all concurrency levels tested (up to 100 concurrent users), `goboxd` maintained a **100% success rate** (200 OK) with zero internal errors or timed-out requests at the default 30s queue timeout setting.

## Raw Benchmark Logs
```text
bash tests/load/load.sh http://localhost:8080
--- Concurrency 1 ---

Summary:
  Total:	3.9239 secs
  Slowest:	0.1048 secs
  Fastest:	0.0147 secs
  Average:	0.0196 secs
  Requests/sec:	50.9701
  
  Total data:	36799 bytes
  Size/request:	183 bytes

Response time histogram:
  0.015 [1]	|
  0.024 [174]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.033 [8]	|■■
  0.042 [10]	|■■
  0.051 [3]	|■
  0.060 [0]	|
  0.069 [1]	|
  0.078 [0]	|
  0.087 [1]	|
  0.096 [0]	|
  0.105 [2]	|


Latency distribution:
  10% in 0.0149 secs
  25% in 0.0152 secs
  50% in 0.0161 secs
  75% in 0.0176 secs
  90% in 0.0295 secs
  95% in 0.0405 secs
  99% in 0.0964 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0000 secs, 0.0016 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0004 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0003 secs
  resp wait:	0.0195 secs, 0.0146 secs, 0.1047 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0003 secs

Status code distribution:
  [200]	200 responses



--- Concurrency 10 ---

Summary:
  Total:	1.3956 secs
  Slowest:	0.1234 secs
  Fastest:	0.0161 secs
  Average:	0.0675 secs
  Requests/sec:	143.3065
  
  Total data:	36800 bytes
  Size/request:	184 bytes

Response time histogram:
  0.016 [1]	|■
  0.027 [1]	|■
  0.038 [0]	|
  0.048 [13]	|■■■■■■■
  0.059 [37]	|■■■■■■■■■■■■■■■■■■■■
  0.070 [74]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.080 [40]	|■■■■■■■■■■■■■■■■■■■■■■
  0.091 [19]	|■■■■■■■■■■
  0.102 [9]	|■■■■■
  0.113 [2]	|■
  0.123 [4]	|■■


Latency distribution:
  10% in 0.0497 secs
  25% in 0.0588 secs
  50% in 0.0651 secs
  75% in 0.0751 secs
  90% in 0.0879 secs
  95% in 0.0978 secs
  99% in 0.1206 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0001 secs, 0.0000 secs, 0.0017 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0015 secs
  req write:	0.0001 secs, 0.0000 secs, 0.0012 secs
  resp wait:	0.0673 secs, 0.0160 secs, 0.1216 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0003 secs

Status code distribution:
  [200]	200 responses



--- Concurrency 50 ---

Summary:
  Total:	1.4481 secs
  Slowest:	0.4425 secs
  Fastest:	0.0327 secs
  Average:	0.3205 secs
  Requests/sec:	138.1120
  
  Total data:	36800 bytes
  Size/request:	184 bytes

Response time histogram:
  0.033 [1]	|
  0.074 [4]	|■■
  0.115 [7]	|■■■
  0.156 [6]	|■■■
  0.197 [7]	|■■■
  0.238 [5]	|■■
  0.279 [7]	|■■■
  0.320 [15]	|■■■■■■■
  0.361 [86]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.402 [48]	|■■■■■■■■■■■■■■■■■■■■■■
  0.442 [14]	|■■■■■■■


Latency distribution:
  10% in 0.1702 secs
  25% in 0.3174 secs
  50% in 0.3438 secs
  75% in 0.3756 secs
  90% in 0.3962 secs
  95% in 0.4093 secs
  99% in 0.4377 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0006 secs, 0.0000 secs, 0.0064 secs
  DNS-lookup:	0.0004 secs, 0.0000 secs, 0.0048 secs
  req write:	0.0001 secs, 0.0000 secs, 0.0010 secs
  resp wait:	0.3198 secs, 0.0316 secs, 0.4424 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0013 secs

Status code distribution:
  [200]	200 responses



--- Concurrency 100 ---

Summary:
  Total:	1.4701 secs
  Slowest:	0.7671 secs
  Fastest:	0.0352 secs
  Average:	0.5588 secs
  Requests/sec:	136.0466
  
  Total data:	36800 bytes
  Size/request:	184 bytes

Response time histogram:
  0.035 [1]	|
  0.108 [8]	|■■■
  0.182 [10]	|■■■■
  0.255 [9]	|■■■
  0.328 [11]	|■■■■
  0.401 [10]	|■■■■
  0.474 [9]	|■■■
  0.548 [12]	|■■■■■
  0.621 [10]	|■■■■
  0.694 [15]	|■■■■■■
  0.767 [105]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■


Latency distribution:
  10% in 0.1973 secs
  25% in 0.4164 secs
  50% in 0.6955 secs
  75% in 0.7115 secs
  90% in 0.7248 secs
  95% in 0.7310 secs
  99% in 0.7642 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0030 secs, 0.0000 secs, 0.0134 secs
  DNS-lookup:	0.0025 secs, 0.0000 secs, 0.0126 secs
  req write:	0.0002 secs, 0.0000 secs, 0.0027 secs
  resp wait:	0.5555 secs, 0.0336 secs, 0.7558 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0072 secs

Status code distribution:
  [200]	200 responses
```

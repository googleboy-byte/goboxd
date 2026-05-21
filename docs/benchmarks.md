# Benchmarks

Benchmarks for `goboxd` at various concurrency levels with the bounded concurrency queue active.

## Environment
- OS: Linux
- Go: 1.23
- Docker: 27.x

## Results (Python 3 - Hello World)

| Concurrency | Requests | Requests/sec | p50 (sec) | p95 (sec) | p99 (sec) |
|-------------|----------|--------------|-----------|-----------|-----------|
| 1           | 200      | 41.13        | 0.0189    | 0.0592    | 0.1383    |
| 10          | 200      | 85.58        | 0.1030    | 0.2152    | 0.2816    |
| 50          | 200      | 101.13       | 0.4268    | 0.6419    | 0.6947    |
| 100         | 200      | 103.13       | 0.8761    | 1.0493    | 1.0807    |

> [!NOTE]
> Bounded concurrency queue is **ACTIVE** (limit: runtime.NumCPU). p99 increases at c=50 and c=100 due to queuing, but server remains stable with 0% error rate.

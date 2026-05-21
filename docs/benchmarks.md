# Benchmarks

Benchmarks for `goboxd` at various concurrency levels.

## Environment
- OS: Linux
- Go: 1.23
- Docker: 27.x

## Results (Python 3 - Hello World)

| Concurrency | Requests | Requests/sec | p50 (sec) | p95 (sec) | p99 (sec) |
|-------------|----------|--------------|-----------|-----------|-----------|
| 1           | 200      | 59.93        | 0.0165    | 0.0183    | 0.0223    |
| 10          | 200      | 151.30       | 0.0586    | 0.1172    | 0.1427    |
| 50          | 200      | 150.12       | 0.1986    | 0.7658    | 0.9329    |
| 100         | 200      | 139.80       | 0.4870    | 1.2559    | 1.3325    |

> [!NOTE]
> No bounded concurrency queue in Stage 1. p99 degradation at c=50 and c=100 is expected. Stage 2 will add queuing.

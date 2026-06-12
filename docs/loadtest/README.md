# Stage 3 Load Test Results

This document contains the container configurations, tool details, results, and replication instructions for the goboxd load test challenge.

## Container Resource Limits
- **CPUs**: 2 vCPUs (`--cpus=2`)
- **Memory**: 2 GB RAM (`--memory=2g`)

## Load Testing Setup
- **Tool**: [Vegeta](https://github.com/tsenart/vegeta)
- **Workload**: `MemoryHog.java` (allocates and touches 150 MB of memory per execution)
- **RPS Ladder**: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400 requests/second
- **Step Duration**: 30 seconds per rate step
- **Request Timeout**: 10 seconds

## Results Summary
- **Breaking-point RPS**: **< 5 RPS** (the service broke at the very first step of 5 RPS, yielding a 90.67% error rate).
- **First Failure Mode**: **Concurrency Limiter / Queue Buildup Timeout**
  - **Explanation**: 
    - The server's `MaxConcurrentJobs` defaults to `runtime.NumCPU()`. Inside the container (limited to 2 CPUs), this value is **2**.
    - Compiling and executing `MemoryHog.java` is extremely expensive: the compilation takes about **1.5s - 1.6s**, and execution sleeps for **1.0s** (total processing time per request is **~2.5s**).
    - Thus, the maximum theoretical throughput of the server is `2 jobs / 2.5 seconds = 0.8 requests/second`.
    - At the lowest step of **5 RPS**, the offered load exceeds the server's processing capacity by more than 6x. The concurrency queue builds up immediately.
    - Since the client-side timeout is set to 10s, requests waiting in the queue for longer than 10 seconds are aborted by the client (Vegeta), leading to connection termination and 100% failure rates at 10 RPS and above.
    - The service **degraded gracefully** without crashing; it did not suffer from memory leaks or Out-Of-Memory (OOM) failures, and returned to full health once the load subsided.

## How to Reproduce

1. **Rebuild the container with resource limits**:
   ```bash
   sudo make run
   ```

2. **Execute the load test script**:
   ```bash
   ./docs/loadtest/load-test.sh
   ```
   This will output a JSON report for each rate step inside `docs/loadtest/` (preserving all stage-by-stage data) and write the summary into `docs/loadtest/results.csv`.

3. **Generate the plots**:
   ```bash
   python3 docs/loadtest/plot.py
   ```
   This will generate the required `breaking-point.png` and `latency.png` plots inside `docs/loadtest/`.

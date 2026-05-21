#!/bin/bash
set -e

# Basic load test script using curl in a loop.
# This serves as a placeholder until k6 or vegeta is integrated in Stage 3.

SERVER_URL=${1:-"http://localhost:8080"}
CONCURRENCY=5
REQUESTS=20

echo "--- Starting Load Test (Concurrency: $CONCURRENCY, Total Requests: $REQUESTS) ---"

run_load() {
    for i in $(seq 1 $REQUESTS); do
        curl -s -X POST -H "Content-Type: application/json" \
            -d '{"language":"py3","source":"print(\"hello\")","tests":[{"stdin":"","expected_stdout":"hello\n"}]}' \
            "$SERVER_URL/run" > /dev/null &
        
        if (( i % CONCURRENCY == 0 )); then
            wait
        fi
    done
    wait
}

time run_load

echo "--- Load test completed ---"

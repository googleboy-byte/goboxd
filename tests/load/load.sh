#!/bin/bash
SERVER_URL=${1:-"http://localhost:8080"}
PAYLOAD='{"language":"py3","source":"print(1)","tests":[{"stdin":"","expected_stdout":"1\n"}]}'

for C in 1 10 50 100; do
    echo "--- Concurrency $C ---"
    hey -n 200 -c $C -m POST \
        -H "Content-Type: application/json" \
        -d "$PAYLOAD" \
        "$SERVER_URL/run"
done

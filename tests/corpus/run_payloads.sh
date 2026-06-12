#!/usr/bin/env bash
# tests/corpus/run_payloads.sh
# Run payloads from the payloads/ directory against the live server.
# Usage: bash tests/corpus/run_payloads.sh [SERVER_URL]

set -euo pipefail

SERVER_URL="${1:-http://localhost:8080}"
PASS=0
FAIL=0

green() { echo -e "\033[32m$1\033[0m"; }
red()   { echo -e "\033[31m$1\033[0m"; }

echo "=============================="
echo " Running Custom Payloads"
echo " Server: $SERVER_URL"
echo "=============================="

# Check server health first
if ! curl -sf "$SERVER_URL/healthz" > /dev/null; then
    echo "Server not ready or unreachable at $SERVER_URL. Please run 'sudo make run' first."
    exit 1
fi

# Find all JSON files in payloads directory
PAYLOAD_FILES=$(find payloads -type f -name "*.json" | sort)

for filepath in $PAYLOAD_FILES; do
    # Extract language from parent directory name
    lang=$(basename "$(dirname "$filepath")")
    # Expected status is the filename without .json
    expected_status=$(basename "$filepath" .json)
    
    echo -n "Testing $lang ($expected_status)... "
    
    # Send request
    resp=$(curl -s -X POST -H "Content-Type: application/json" -d @"$filepath" "$SERVER_URL/run")
    
    # Parse status
    status=$(echo "$resp" | jq -r '.status // "null"')
    
    if [ "$status" == "$expected_status" ]; then
        green "✅ PASSED"
        echo "  Expected: $expected_status"
        echo "  Got:      $status"
        echo "  Response: $resp"
        ((PASS++)) || true
    else
        red "❌ FAILED"
        echo ""
        echo "  Expected: $expected_status"
        echo "  Got:      $status"
        echo "  Response: $resp"
        ((FAIL++)) || true
    fi
done

echo ""
echo "=============================="
TOTAL=$((PASS + FAIL))
echo " Results: $PASS/$TOTAL passed"
if [ "$FAIL" -eq 0 ]; then
    green " ALL PAYLOADS PASSED"
    exit 0
else
    red " $FAIL PAYLOAD(S) FAILED"
    exit 1
fi

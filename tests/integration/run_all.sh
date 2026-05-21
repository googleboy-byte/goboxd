#!/bin/bash
set -e

# This script runs a basic smoke test for all registered languages.
# It assumes goboxd is running at localhost:8080.

SERVER_URL=${1:-"http://localhost:8080"}

test_lang() {
    local lang=$1
    local source=$2
    local expected_status=$3
    
    echo "Testing $lang..."
    
    # Use jq -n to build the JSON properly with real newlines
    local payload=$(jq -n \
        --arg lang "$lang" \
        --arg src "$source" \
        '{language: $lang, source: $src, tests: [{stdin: "", expected_stdout: "hello\n"}]}')

    local resp=$(curl -s -X POST -H "Content-Type: application/json" -d "$payload" "$SERVER_URL/run")
    
    local status=$(echo "$resp" | jq -r '.status')
    
    if [ "$status" == "$expected_status" ]; then
        echo "✅ $lang: $status"
    else
        echo "❌ $lang: Expected $expected_status, got $status"
        echo "Full response: $resp"
        exit 1
    fi
}

echo "--- Integration Tests ---"

# 1. Python 3
test_lang "py3" "print('hello')" "accepted"

# 2. C++
test_lang "cpp" $'#include <iostream>\nint main() { std::cout << "hello" << std::endl; return 0; }' "accepted"

# 3. Bash
test_lang "bash" "echo hello" "accepted"

# 4. Rust
test_lang "rust" $'fn main() { println!("hello"); }' "accepted"

echo "--- All integration tests passed! ---"

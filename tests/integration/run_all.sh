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

# 5. Java (requires source_filename and artifact_filename)
echo "Testing java..."
JAVA_RESP=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"java",
  "source":"public class Hello { public static void main(String[] args) { System.out.println(\"hello\"); } }",
  "source_filename": "Hello.java",
  "artifact_filename": "Hello",
  "tests":[{"stdin":"","expected_stdout":"hello\n"}]
}' "$SERVER_URL/run")
JAVA_STATUS=$(echo "$JAVA_RESP" | jq -r '.status')
if [ "$JAVA_STATUS" == "accepted" ]; then
    echo "✅ java: $JAVA_STATUS"
else
    echo "❌ java: Expected accepted, got $JAVA_STATUS"
    echo "Full response: $JAVA_RESP"
    exit 1
fi

# 6. C
test_lang "c" $'#include <stdio.h>\nint main() { printf("hello\\n"); return 0; }' "accepted"

# 7. JavaScript
test_lang "js" "console.log('hello')" "accepted"

echo "--- All integration tests passed! ---"

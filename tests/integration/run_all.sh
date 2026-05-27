#!/bin/bash
set -e

# This script runs a basic smoke test for all registered languages.
# It assumes goboxd is running at localhost:8080.

SERVER_URL=${1:-"http://localhost:8080"}

# Fetch registered languages from /info
echo "Discovering registered languages..."
# Get IDs, convert to single line space-separated
REGISTERED_LANGS=$(curl -s "$SERVER_URL/info" | jq -r '.languages[].id' | xargs)

test_lang() {
    local lang=$1
    local source=$2
    local expected_status=$3
    
    # Check if language is registered (match full word)
    if [[ ! " $REGISTERED_LANGS " =~ " $lang " ]]; then
        echo "[-] $lang: Skipped (not registered)"
        return 0
    fi

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
if [[ " $REGISTERED_LANGS " =~ " java " ]]; then
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
else
    echo "⏭️  java: Skipped (not registered)"
fi

# 6. C
test_lang "c" $'#include <stdio.h>\nint main() { printf("hello\\n"); return 0; }' "accepted"

# 7. JavaScript
test_lang "js" "console.log('hello')" "accepted"

# 8. Verilog
test_lang "verilog" "module main; initial begin \$display(\"hello\"); \$finish; end endmodule" "accepted"

# 9. Go
test_lang "go" "package main; import \"fmt\"; func main() { fmt.Println(\"hello\") }" "accepted"

# 10. Kotlin
test_lang "kotlin" "fun main() { println(\"hello\") }" "accepted"

# 11. C# (Mono)
test_lang "csharp" "using System; class Hello { static void Main() { Console.WriteLine(\"hello\"); } }" "accepted"

# 12. Ruby
test_lang "ruby" "puts 'hello'" "accepted"

# 13. Lua
test_lang "lua" "print('hello')" "accepted"

# 14. OCaml
test_lang "ocaml" "print_endline \"hello\"" "accepted"

# 15. Swift
test_lang "swift" "print(\"hello\")" "accepted"

# 16. Zig
test_lang "zig" 'const std = @import("std"); pub fn main() !void { try std.io.getStdOut().writer().writeAll("hello\n"); }' "accepted"

echo "--- Integration tests completed! ---"

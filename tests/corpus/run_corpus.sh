#!/usr/bin/env bash
# tests/corpus/run_corpus.sh
# Corpus test suite for goboxd. Runs against a live server.
# Usage: bash tests/corpus/run_corpus.sh [SERVER_URL]
# Default SERVER_URL: http://localhost:8080

set -euo pipefail

SERVER_URL="${1:-http://localhost:8080}"
PASS=0
FAIL=0
SKIP=0

green() { echo -e "\033[32m✅ $1\033[0m"; }
red()   { echo -e "\033[31m❌ $1\033[0m"; }

pass() { green "$1"; ((PASS++)) || true; }
fail() { red   "$1"; ((FAIL++)) || true; }
skip() { echo -e "[-] $1: Skipped (not registered)"; ((SKIP++)) || true; }

# assert_status <test_name> <expected_top_status> <json_response>
assert_status() {
    local name="$1" expected="$2" resp="$3"
    local got
    got=$(echo "$resp" | jq -r '.status // "null"')
    if [ "$got" == "$expected" ]; then
        pass "$name"
    else
        fail "$name — expected '$expected', got '$got'"
        echo "  Response: $resp"
    fi
}

# assert_error_code <test_name> <expected_code> <json_response>
assert_error_code() {
    local name="$1" expected="$2" resp="$3"
    local got
    got=$(echo "$resp" | jq -r '.error.code // "null"')
    if [ "$got" == "$expected" ]; then
        pass "$name"
    else
        fail "$name — expected error.code '$expected', got '$got'"
        echo "  Response: $resp"
    fi
}

# assert_http <test_name> <expected_http_code> <url> [curl_args...]
assert_http() {
    local name="$1" expected="$2" url="$3"
    shift 3
    local got
    got=$(curl -s -o /dev/null -w "%{http_code}" "$@" "$url")
    if [ "$got" == "$expected" ]; then
        pass "$name"
    else
        fail "$name — expected HTTP $expected, got $got"
    fi
}

# post_run <json_body> — returns response body
post_run() {
    local tmp=$(mktemp)
    echo "$1" > "$tmp"
    curl -s -X POST -H "Content-Type: application/json" -d @"$tmp" "$SERVER_URL/run"
    rm "$tmp"
}

echo "=============================="
echo " goboxd corpus test suite"
echo " Server: $SERVER_URL"
echo "=============================="

# ── Wait for server ────────────────────────────────────────────────────────────
echo ""
echo "── Waiting for server ──"
for i in $(seq 1 20); do
    if curl -sf "$SERVER_URL/healthz" > /dev/null 2>&1; then break; fi
    sleep 1
done
curl -sf "$SERVER_URL/healthz" > /dev/null || { echo "Server not ready. Aborting."; exit 1; }
pass "server reachable"

# ── Health endpoints ───────────────────────────────────────────────────────────
echo ""
echo "── Health endpoints ──"

HEALTHZ=$(curl -s "$SERVER_URL/healthz")
HEALTHZ_STATUS=$(echo "$HEALTHZ" | jq -r '.status')
[ "$HEALTHZ_STATUS" == "ok" ] && pass "/healthz returns ok" || fail "/healthz status: $HEALTHZ_STATUS"

READYZ=$(curl -s "$SERVER_URL/readyz")
READYZ_STATUS=$(echo "$READYZ" | jq -r '.status')
[ "$READYZ_STATUS" == "ok" ] && pass "/readyz returns ok" || fail "/readyz status: $READYZ_STATUS"

INFO=$(curl -s "$SERVER_URL/info")
echo "$INFO" | jq -e '.build_info.version' > /dev/null && pass "/info has build_info.version" || fail "/info missing build_info.version"
echo "$INFO" | jq -e '.languages | length > 0' > /dev/null && pass "/info has languages" || fail "/info missing languages"
echo "$INFO" | jq -e '.limits.max_source_bytes' > /dev/null && pass "/info has limits.max_source_bytes" || fail "/info missing limits.max_source_bytes"
echo "$INFO" | jq -e '.stats.jobs_total' > /dev/null && pass "/info has stats.jobs_total" || fail "/info missing stats.jobs_total"

# Discover registered languages
REGISTERED_LANGS=$(echo "$INFO" | jq -r '.languages[].id' | xargs)

# ── Happy path: hello world per language ──────────────────────────────────────
echo ""
echo "── Happy path: hello world ──"

if [[ " $REGISTERED_LANGS " =~ " py3 " ]]; then
    R=$(post_run '{"language":"py3","source":"print(\"hello\")","tests":[{"stdin":"","expected_stdout":"hello\n"}]}')
    assert_status "py3 hello world" "accepted" "$R"
else skip "py3"; fi

if [[ " $REGISTERED_LANGS " =~ " cpp " ]]; then
    R=$(post_run '{"language":"cpp","source":"#include<iostream>\nint main(){std::cout<<\"hello\"<<std::endl;}","tests":[{"stdin":"","expected_stdout":"hello\n"}]}')
    assert_status "cpp hello world" "accepted" "$R"
else skip "cpp"; fi

if [[ " $REGISTERED_LANGS " =~ " c " ]]; then
    R=$(post_run '{"language":"c","source":"#include<stdio.h>\nint main(){printf(\"hello\\n\");}","tests":[{"stdin":"","expected_stdout":"hello\n"}]}')
    assert_status "c hello world" "accepted" "$R"
else skip "c"; fi

if [[ " $REGISTERED_LANGS " =~ " bash " ]]; then
    R=$(post_run '{"language":"bash","source":"echo hello","tests":[{"stdin":"","expected_stdout":"hello\n"}]}')
    assert_status "bash hello world" "accepted" "$R"
else skip "bash"; fi

if [[ " $REGISTERED_LANGS " =~ " js " ]]; then
    R=$(post_run '{"language":"js","source":"console.log(\"hello\")","tests":[{"stdin":"","expected_stdout":"hello\n"}]}')
    assert_status "js hello world" "accepted" "$R"
else skip "js"; fi

if [[ " $REGISTERED_LANGS " =~ " rust " ]]; then
    R=$(post_run '{"language":"rust","source":"fn main(){println!(\"hello\");}","tests":[{"stdin":"","expected_stdout":"hello\n"}]}')
    assert_status "rust hello world" "accepted" "$R"
else skip "rust"; fi

if [[ " $REGISTERED_LANGS " =~ " java " ]]; then
    R=$(post_run '{
      "language":"java",
      "source":"public class Hello { public static void main(String[] a) { System.out.println(\"hello\"); } }",
      "source_filename":"Hello.java",
      "artifact_filename":"Hello",
      "tests":[{"stdin":"","expected_stdout":"hello\n"}]
    }')
    assert_status "java hello world" "accepted" "$R"
else skip "java"; fi

if [[ " $REGISTERED_LANGS " =~ " verilog " ]]; then
    R=$(post_run '{"language":"verilog","source":"module main; initial begin $display(\"hello\"); $finish; end endmodule","tests":[{"stdin":"","expected_stdout":"hello\n"}]}')
    assert_status "verilog hello world" "accepted" "$R"
else skip "verilog"; fi

# ── stdin echo ─────────────────────────────────────────────────────────────────
echo ""
echo "── stdin echo ──"

R=$(post_run '{"language":"py3","source":"import sys; print(sys.stdin.read().strip())","tests":[{"stdin":"hello","expected_stdout":"hello\n"}]}')
assert_status "py3 stdin echo" "accepted" "$R"

R=$(post_run '{"language":"c","source":"#include<stdio.h>\nint main(){char b[64];fgets(b,64,stdin);printf(\"%s\",b);}","tests":[{"stdin":"hello\n","expected_stdout":"hello\n"}]}')
assert_status "c stdin echo" "accepted" "$R"

# ── multiple test cases ────────────────────────────────────────────────────────
echo ""
echo "── multiple test cases ──"

R=$(post_run '{
  "language":"py3",
  "source":"x=int(input()); print(x*2)",
  "tests":[
    {"stdin":"1\n","expected_stdout":"2\n"},
    {"stdin":"5\n","expected_stdout":"10\n"},
    {"stdin":"0\n","expected_stdout":"0\n"}
  ]
}')
assert_status "py3 multi-test accepted" "accepted" "$R"
TEST_COUNT=$(echo "$R" | jq '.tests | length')
[ "$TEST_COUNT" == "3" ] && pass "py3 multi-test: 3 results returned" || fail "py3 multi-test: expected 3 results, got $TEST_COUNT"

# ── wrong output ───────────────────────────────────────────────────────────────
echo ""
echo "── status: wrong_output ──"

R=$(post_run '{"language":"py3","source":"print(\"wrong\")","tests":[{"stdin":"","expected_stdout":"right\n"}]}')
assert_status "py3 wrong_output" "wrong_output" "$R"
BUILD_STATUS=$(echo "$R" | jq -r '.build.status')
[ "$BUILD_STATUS" == "ok" ] && pass "py3 wrong_output: build.status ok" || fail "py3 wrong_output: build.status was $BUILD_STATUS"

# ── whitespace mismatch ────────────────────────────────────────────────────────
echo ""
echo "── status: output_whitespace_mismatch ──"

R=$(post_run '{"language":"py3","source":"print(\"hello\")","tests":[{"stdin":"","expected_stdout":"hello"}]}')
GOT=$(echo "$R" | jq -r '.status')
# acceptable: output_whitespace_mismatch or wrong_output depending on impl
if [ "$GOT" == "output_whitespace_mismatch" ] || [ "$GOT" == "wrong_output" ]; then
    pass "py3 whitespace mismatch detected ($GOT)"
else
    fail "py3 whitespace mismatch — expected mismatch status, got '$GOT'"
fi

# ── not_executed after build failure ──────────────────────────────────────────
echo ""
echo "── build failure → not_executed ──"

R=$(post_run '{"language":"cpp","source":"this is not valid c++","tests":[{"stdin":"","expected_stdout":"hello\n"},{"stdin":"","expected_stdout":"world\n"}]}')
assert_status "cpp build_failed top-level" "build_failed" "$R"
BUILD_STATUS=$(echo "$R" | jq -r '.build.status')
[ "$BUILD_STATUS" == "failed" ] && pass "cpp build.status failed" || fail "cpp build.status: expected failed, got $BUILD_STATUS"
ALL_NOT_EXECUTED=$(echo "$R" | jq '[.tests[].status == "not_executed"] | all')
[ "$ALL_NOT_EXECUTED" == "true" ] && pass "cpp all tests not_executed after build failure" || fail "cpp some tests not not_executed after build failure"

R=$(post_run '{"language":"py3","source":"def broken(","tests":[{"stdin":"","expected_stdout":"x\n"}]}')
GOT=$(echo "$R" | jq -r '.status')
if [ "$GOT" == "build_failed" ] || [ "$GOT" == "runtime_error" ]; then
    pass "py3 syntax error detected ($GOT)"
else
    fail "py3 syntax error — expected build_failed or runtime_error, got '$GOT'"
fi

# ── flag override ──────────────────────────────────────────────────────────────
echo ""
echo "── flag override ──"

R=$(post_run '{
  "language":"cpp",
  "source":"#include<iostream>\nint main(){std::cout<<\"hi\"<<std::endl;}",
  "build":{"flags":["-O2","-Wall"]},
  "tests":[{"stdin":"","expected_stdout":"hi\n"}]
}')
assert_status "cpp with allowed flags" "accepted" "$R"

# ── runtime error ──────────────────────────────────────────────────────────────
echo ""
echo "── status: runtime_error ──"

R=$(post_run '{"language":"c","source":"#include<stdlib.h>\nint main(){abort();}","tests":[{"stdin":"","expected_stdout":""}]}')
GOT=$(echo "$R" | jq -r '.status')
[ "$GOT" == "runtime_error" ] && pass "c abort() → runtime_error" || fail "c abort() — expected runtime_error, got '$GOT'"

R=$(post_run '{"language":"py3","source":"raise RuntimeError(\"boom\")","tests":[{"stdin":"","expected_stdout":""}]}')
GOT=$(echo "$R" | jq -r '.tests[0].status')
[ "$GOT" == "runtime_error" ] && pass "py3 exception → runtime_error" || fail "py3 exception — expected runtime_error, got '$GOT'"

# ── timeout ───────────────────────────────────────────────────────────────────
echo ""
echo "── status: time_exceeded ──"

R=$(post_run '{"language":"py3","source":"while True: pass","tests":[{"stdin":"","expected_stdout":""}]}')
assert_status "py3 infinite loop → time_exceeded" "time_exceeded" "$R"

R=$(post_run '{"language":"c","source":"#include<stdio.h>\nint main(){for(;;);}","tests":[{"stdin":"","expected_stdout":""}]}')
assert_status "c infinite loop → time_exceeded" "time_exceeded" "$R"

# ── empty output ──────────────────────────────────────────────────────────────
echo ""
echo "── edge: empty output ──"

R=$(post_run '{"language":"py3","source":"pass","tests":[{"stdin":"","expected_stdout":""}]}')
assert_status "py3 empty output accepted" "accepted" "$R"

# ── 50 test cases (max) ───────────────────────────────────────────────────────
echo ""
echo "── edge: 50 test cases ──"

TESTS_50=$(python3 -c "
import json
tests = [{'stdin': str(i)+'\n', 'expected_stdout': str(i*2)+'\n'} for i in range(50)]
print(json.dumps(tests))
")
R=$(post_run "{\"language\":\"py3\",\"source\":\"x=int(input()); print(x*2)\",\"tests\":$TESTS_50}")
assert_status "py3 50 tests accepted" "accepted" "$R"
COUNT=$(echo "$R" | jq '.tests | length')
[ "$COUNT" == "50" ] && pass "py3 50 tests: 50 results returned" || fail "py3 50 tests: got $COUNT results"

# ── source at size boundary ───────────────────────────────────────────────────
echo ""
echo "── edge: source size boundary ──"

# 256KiB - 1 byte: should succeed (padded with comments)
BIG_SOURCE_FILE=$(mktemp)
python3 -c "
pad = '#' * (262143 - len('print(\"hi\")') - 1)
print('print(\"hi\") ' + pad)
" > "$BIG_SOURCE_FILE"
PAYLOAD_FILE=$(mktemp)
python3 -c "
import json
with open('$BIG_SOURCE_FILE', 'r') as f:
    src = f.read()
print(json.dumps({'language':'py3','source':src,'tests':[{'stdin':'','expected_stdout':'hi\n'}]}))
" > "$PAYLOAD_FILE"
R=$(curl -s -X POST -H "Content-Type: application/json" -d @"$PAYLOAD_FILE" "$SERVER_URL/run")
assert_status "py3 source at 256KiB-1 accepted" "accepted" "$R"
rm "$BIG_SOURCE_FILE" "$PAYLOAD_FILE"

# ── adversarial: path traversal ───────────────────────────────────────────────
echo ""
echo "── adversarial: path traversal ──"

R=$(post_run '{"language":"cpp","source":"int main(){}","source_filename":"../../etc/passwd","tests":[{"stdin":"","expected_stdout":""}]}')
assert_error_code "path traversal source_filename rejected" "invalid_filename" "$R"

R=$(post_run '{"language":"cpp","source":"int main(){}","artifact_filename":"../evil","tests":[{"stdin":"","expected_stdout":""}]}')
assert_error_code "path traversal artifact_filename rejected" "invalid_filename" "$R"

# ── adversarial: disallowed flags ─────────────────────────────────────────────
echo ""
echo "── adversarial: disallowed flags ──"

R=$(post_run '{"language":"cpp","source":"int main(){}","build":{"flags":["-fplugin=evil.so"]},"tests":[{"stdin":"","expected_stdout":""}]}')
assert_error_code "cpp disallowed build flag rejected" "disallowed_flag" "$R"

R=$(post_run '{"language":"cpp","source":"int main(){}","build":{"flags":["--specs=/evil"]},"tests":[{"stdin":"","expected_stdout":""}]}')
assert_error_code "cpp --specs flag rejected" "disallowed_flag" "$R"

# ── adversarial: oversize body ────────────────────────────────────────────────
echo ""
echo "── adversarial: oversize body ──"

OVERSIZE_FILE=$(mktemp)
python3 -c "print('A' * 300000)" > "$OVERSIZE_FILE"
PAYLOAD_FILE=$(mktemp)
python3 -c "
import json
with open('$OVERSIZE_FILE', 'r') as f:
    src = f.read()
print(json.dumps({'language':'py3','source':src,'tests':[{'stdin':'','expected_stdout':''}]}))
" > "$PAYLOAD_FILE"
R=$(curl -s -X POST -H "Content-Type: application/json" -d @"$PAYLOAD_FILE" "$SERVER_URL/run")
GOT_CODE=$(echo "$R" | jq -r '.error.code // .status // "null"')
[ "$GOT_CODE" != "accepted" ] && pass "oversize body rejected" || fail "oversize body was accepted"
rm "$OVERSIZE_FILE" "$PAYLOAD_FILE"

# ── adversarial: unknown language ─────────────────────────────────────────────
echo ""
echo "── adversarial: unknown language ──"

R=$(post_run '{"language":"cobol","source":"x","tests":[{"stdin":"","expected_stdout":""}]}')
assert_error_code "unknown language rejected" "unknown_language" "$R"

# ── adversarial: malformed JSON ───────────────────────────────────────────────
echo ""
echo "── adversarial: malformed JSON ──"

R=$(curl -s -X POST -H "Content-Type: application/json" -d '{bad json' "$SERVER_URL/run")
assert_error_code "malformed JSON rejected" "invalid_json" "$R"

# ── adversarial: missing required fields ──────────────────────────────────────
echo ""
echo "── adversarial: missing fields ──"

R=$(post_run '{"language":"py3","tests":[{"stdin":"","expected_stdout":""}]}')
assert_error_code "missing source rejected" "bad_request" "$R"

R=$(post_run '{"source":"print(1)","tests":[{"stdin":"","expected_stdout":"1\n"}]}')
assert_error_code "missing language rejected" "bad_request" "$R"

R=$(post_run '{"language":"py3","source":"print(1)"}')
assert_error_code "missing tests rejected" "bad_request" "$R"

R=$(post_run '{"language":"py3","source":"print(1)","tests":[]}')
assert_error_code "empty tests array rejected" "bad_request" "$R"

# ── adversarial: Java missing filenames ───────────────────────────────────────
echo ""
echo "── adversarial: Java without required filenames ──"

R=$(post_run '{"language":"java","source":"public class X{}","tests":[{"stdin":"","expected_stdout":""}]}')
assert_error_code "java without source_filename rejected" "bad_request" "$R"

R=$(post_run '{"language":"java","source":"public class X{}","source_filename":"X.java","tests":[{"stdin":"","expected_stdout":""}]}')
assert_error_code "java without artifact_filename rejected" "bad_request" "$R"

# ── adversarial: too many tests ───────────────────────────────────────────────
echo ""
echo "── adversarial: too many tests ──"

TESTS_51=$(python3 -c "
import json
tests = [{'stdin': '', 'expected_stdout': ''} for _ in range(51)]
print(json.dumps(tests))
")
R=$(post_run "{\"language\":\"py3\",\"source\":\"pass\",\"tests\":$TESTS_51}")
assert_error_code "51 tests rejected" "bad_request" "$R"

# ── load: 200 requests at c=50 ────────────────────────────────────────────────
echo ""
echo "── load: 200 requests at c=50 ──"

if command -v hey &> /dev/null; then
    LOAD_RESULT=$(hey -n 200 -c 50 -m POST \
        -H "Content-Type: application/json" \
        -d '{"language":"py3","source":"print(\"hi\")","tests":[{"stdin":"","expected_stdout":"hi\n"}]}' \
        "$SERVER_URL/run" 2>&1)
    SUCCESS=$(echo "$LOAD_RESULT" | grep "\[200\]" | awk '{print $2}')
    if [ "$SUCCESS" == "200" ]; then
        pass "load: all 200 requests returned 200"
    else
        fail "load: only $SUCCESS/200 requests returned 200"
        echo "$LOAD_RESULT" | grep "Status code"
    fi
else
    echo "  (skipped — hey not found)"
fi

# ── summary ───────────────────────────────────────────────────────────────────
echo ""
echo "=============================="
TOTAL=$((PASS + FAIL))
echo " Results: $PASS/$TOTAL passed ($SKIP skipped)"
if [ "$FAIL" -eq 0 ]; then
    green " ALL CORPUS TESTS PASSED"
    exit 0
else
    red " $FAIL TEST(S) FAILED"
    exit 1
fi

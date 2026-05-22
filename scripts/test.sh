#!/bin/bash
set -e

SERVER=${1:-"http://localhost:8080"}

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1"; echo "Response: $2"; exit 1; }

check_field() {
    local desc=$1
    local expected=$2
    local actual=$3
    if [ "$actual" = "$expected" ]; then
        pass "$desc"
    else
        fail "$desc (expected $expected got $actual)" "$actual"
    fi
}

echo "--- Waiting for server ---"
for i in $(seq 1 10); do
    curl -s -o /dev/null --connect-timeout 1 $SERVER/healthz && break
    echo "Waiting..."
    sleep 2
done

echo ""
echo "--- Health endpoints ---"

R=$(curl -s $SERVER/healthz)
check_field "healthz status" "ok" $(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['status'])")

R=$(curl -s $SERVER/readyz)
check_field "readyz status" "ok" $(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['status'])")

R=$(curl -s $SERVER/info)
V=$(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['build_info']['version'])")
[ -n "$V" ] && pass "info build_info.version present" || fail "info missing version" "$R"

echo ""
echo "--- POST /run ---"

# py3 accepted
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"py3",
  "source":"print(\"hello\")",
  "tests":[{"stdin":"","expected_stdout":"hello\n"}]
}' $SERVER/run)
check_field "py3 accepted" "accepted" $(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['status'])")

# cpp accepted
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"cpp",
  "source":"#include<iostream>\nint main(){std::cout<<\"hello\";return 0;}",
  "tests":[{"stdin":"","expected_stdout":"hello"}]
}' $SERVER/run)
check_field "cpp accepted" "accepted" $(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['status'])")

# cpp build failed
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"cpp",
  "source":"not valid c++",
  "tests":[{"stdin":"","expected_stdout":""}]
}' $SERVER/run)
check_field "cpp build_failed" "build_failed" $(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['status'])")

# bash accepted
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"bash",
  "source":"echo hello",
  "tests":[{"stdin":"","expected_stdout":"hello\n"}]
}' $SERVER/run)
check_field "bash accepted" "accepted" $(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['status'])")

# rust accepted
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"rust",
  "source":"fn main() { println!(\"hello\"); }",
  "tests":[{"stdin":"","expected_stdout":"hello\n"}]
}' $SERVER/run)
check_field "rust accepted" "accepted" $(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['status'])")

echo ""
echo "--- Error paths ---"

# unknown language
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"fortran",
  "source":"print *",
  "tests":[{"stdin":"","expected_stdout":""}]
}' $SERVER/run)
CODE=$(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['error']['code'])")
check_field "unknown language 400" "unknown_language" "$CODE"

# disallowed flag
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"cpp",
  "source":"int main(){}",
  "build":{"flags":["-fplugin=evil.so"]},
  "tests":[{"stdin":"","expected_stdout":""}]
}' $SERVER/run)
CODE=$(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['error']['code'])")
check_field "disallowed flag 400" "disallowed_flag" "$CODE"

# path traversal
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"cpp",
  "source":"int main(){}",
  "source_filename":"../../etc/passwd",
  "tests":[{"stdin":"","expected_stdout":""}]
}' $SERVER/run)
CODE=$(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['error']['code'])")
check_field "path traversal 400" "invalid_filename" "$CODE"

echo ""
echo "--- Security ---"

# output truncation
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"py3",
  "source":"print(\"A\"*1024*1024,end=\"\")",
  "tests":[{"stdin":"","expected_stdout":""}]
}' $SERVER/run)
TRUNCATED=$(echo $R | python3 -c "import sys,json; r=json.load(sys.stdin); print('[TRUNCATED]' in r['tests'][0]['stdout'])")
check_field "output truncation marker" "True" "$TRUNCATED"

# per-request limit override
R=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"cpp",
  "source":"#include<iostream>\nint main(){std::cout<<\"hi\";return 0;}",
  "build":{"limits":{"wall_time_s":5,"memory_kb":1048576,"max_processes":100},"flags":["-O2"]},
  "run":{"limits":{"wall_time_s":3,"memory_kb":524288,"max_processes":64}},
  "tests":[{"stdin":"","expected_stdout":"hi"}]
}' $SERVER/run)
check_field "per-request limit override" "accepted" $(echo $R | python3 -c "import sys,json; print(json.load(sys.stdin)['status'])")

echo ""
echo "--- ALL TESTS PASSED ---"

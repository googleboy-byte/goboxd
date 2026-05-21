#!/bin/bash
# Security verification script for goboxd

SERVER_URL=${1:-"http://localhost:8080"}
echo "Verifying security holes for $SERVER_URL..."

# Hole 1: Path Traversal
echo "[Hole 1] Path Traversal via SourceFilename..."
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "Content-Type: application/json" -d '{
  "language":"py3",
  "source":"print(1)",
  "source_filename":"../../etc/passwd",
  "tests":[{"stdin":"","expected_stdout":"1\n"}]
}' $SERVER_URL/run)
if [ "$STATUS" -eq 400 ]; then
  echo "  PASS: Rejected malicious SourceFilename"
else
  echo "  FAIL: Received $STATUS for malicious SourceFilename"
  exit 1
fi

echo "[Hole 1] Path Traversal via ArtifactFilename..."
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "Content-Type: application/json" -d '{
  "language":"py3",
  "source":"print(1)",
  "artifact_filename":"../../etc/passwd",
  "tests":[{"stdin":"","expected_stdout":"1\n"}]
}' $SERVER_URL/run)
if [ "$STATUS" -eq 400 ]; then
  echo "  PASS: Rejected malicious ArtifactFilename"
else
  echo "  FAIL: Received $STATUS for malicious ArtifactFilename"
  exit 1
fi

# Hole 3: Compiler Flag Injection
echo "[Hole 3] Compiler Flag Injection..."
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "Content-Type: application/json" -d '{
  "language":"cpp",
  "source":"int main(){}",
  "build":{"flags":["-fplugin=evil.so"]},
  "tests":[{"stdin":"","expected_stdout":""}]
}' $SERVER_URL/run)
if [ "$STATUS" -eq 400 ]; then
  echo "  PASS: Rejected disallowed build flag"
else
  echo "  FAIL: Received $STATUS for disallowed build flag"
  exit 1
fi

echo "[Hole 3] Run Flag Injection..."
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST -H "Content-Type: application/json" -d '{
  "language":"py3",
  "source":"print(1)",
  "run":{"flags":["-E"]},
  "tests":[{"stdin":"","expected_stdout":"1\n"}]
}' $SERVER_URL/run)
if [ "$STATUS" -eq 400 ]; then
  echo "  PASS: Rejected disallowed run flag"
else
  echo "  FAIL: Received $STATUS for disallowed run flag"
  exit 1
fi

# Hole 4: Request Size Limits
echo "[Hole 4] Request Size Limits..."
# Large source
STATUS=$(python3 -c "import requests; print(requests.post('$SERVER_URL/run', json={'language':'py3','source':'x'*300*1024,'tests':[{'stdin':'','expected_stdout':''}]}).status_code)")
if [ "$STATUS" -eq 400 ]; then
  echo "  PASS: Rejected large source"
else
  echo "  FAIL: Received $STATUS for large source"
  exit 1
fi

# Large stdin
STATUS=$(python3 -c "import requests; print(requests.post('$SERVER_URL/run', json={'language':'py3','source':'print(1)','tests':[{'stdin':'x'*65537,'expected_stdout':'1\n'}]}).status_code)")
if [ "$STATUS" -eq 400 ]; then
  echo "  PASS: Rejected large stdin"
else
  echo "  FAIL: Received $STATUS for large stdin"
  exit 1
fi

# Hole 6: Output Truncation
echo "[Hole 6] Output Truncation Marker..."
OUT=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"py3",
  "source":"print(\"A\"*1024*1024)",
  "tests":[{"stdin":"","expected_stdout":""}]
}' $SERVER_URL/run | jq -r '.tests[0].stdout')
if [[ "$OUT" == *"[TRUNCATED]"* ]]; then
  echo "  PASS: Truncation marker present"
else
  echo "  FAIL: Truncation marker missing"
  exit 1
fi

echo "ALL SECURITY TESTS PASSED"

#!/bin/bash
# Security verification script for goboxd

SERVER_URL=${1:-"http://localhost:8080"}
echo "Waiting for $SERVER_URL to be ready..."
for i in {1..10}; do
  if curl -s $SERVER_URL/healthz > /dev/null; then
    echo "Server is ready!"
    break
  fi
  echo "Waiting..."
  sleep 1
done

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

# Hole 7: Network Isolation
echo "[Hole 7] Network Isolation..."
# Try to connect to 1.1.1.1:80 (external) with a 1s timeout
STATUS_JSON=$(curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"py3",
  "source":"import socket\ntry:\n  socket.create_connection((\"1.1.1.1\", 80), timeout=1)\n  print(\"CONNECTED\")\nexcept Exception as e:\n  print(\"ISOLATED\")",
  "tests":[{"stdin":"","expected_stdout":"ISOLATED\n"}]
}' $SERVER_URL/run)
STATUS=$(echo "$STATUS_JSON" | jq -r '.status')
STDOUT=$(echo "$STATUS_JSON" | jq -r '.tests[0].stdout')

if [ "$STATUS" == "accepted" ] && [ "$STDOUT" == "ISOLATED" ]; then
  echo "  PASS: Network is isolated"
else
  echo "  FAIL: Network is NOT isolated (Status: $STATUS, Output: $STDOUT)"
  exit 1
fi

echo "ALL SECURITY TESTS PASSED"

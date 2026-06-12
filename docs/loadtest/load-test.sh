#!/usr/bin/env bash
# docs/loadtest/load-test.sh

set -euo pipefail

VEGETA=~/go/bin/vegeta
TARGET_FILE="docs/loadtest/target.txt"
BASE_DIR="docs/loadtest"
CHANGES_FILE="${BASE_DIR}/changes.txt"

# Auto-detect next run number
RUN_NUM=1
while [ -d "${BASE_DIR}/run_${RUN_NUM}" ]; do
  RUN_NUM=$((RUN_NUM + 1))
done

RUN_DIR="${BASE_DIR}/run_${RUN_NUM}"
mkdir -p "$RUN_DIR"

RESULTS_CSV="${RUN_DIR}/results.csv"

# Generate metadata file
METADATA="${RUN_DIR}/metadata.md"
CHANGES_CONTENT="Baseline — no changes from default configuration."
if [ -f "$CHANGES_FILE" ]; then
  FILE_CONTENT=$(cat "$CHANGES_FILE")
  # Check if it's not just the empty template
  if [ -n "$FILE_CONTENT" ] && ! echo "$FILE_CONTENT" | grep -q "^No changes. Edit this file"; then
    CHANGES_CONTENT="$FILE_CONTENT"
  fi
fi

cat > "$METADATA" <<EOF
# Run ${RUN_NUM} — $(date '+%Y-%m-%d %H:%M:%S %Z')

## Changes before this run

${CHANGES_CONTENT}

## Container limits

- CPUs: 2
- Memory: 2G

## Test parameters

- RPS ladder: 5, 10, 25, 50, 75, 100, 150, 200, 300, 400
- Duration per step: 30s
- Client timeout: 10s
EOF

# Reset changes file to empty template
cat > "$CHANGES_FILE" <<'TEMPLATE'
No changes. Edit this file before running load-test.sh to document what was changed.
TEMPLATE

echo "==============================="
echo " Load Test Run #${RUN_NUM}"
echo " Output: ${RUN_DIR}/"
echo "==============================="

echo "target_rps,throughput_rps,duration_s,requests,success,failed,error_pct,p50_ms,p95_ms,p99_ms,max_ms" > "$RESULTS_CSV"

for rate in 5 10 25 50 75 100 150 200 300 400; do
  echo "Attacking at rate: ${rate} rps..."
  report_json="${RUN_DIR}/report-${rate}.json"

  $VEGETA attack -rate=${rate}/1s -duration=30s -timeout=10s -targets="$TARGET_FILE" \
    | $VEGETA report -type=json > "$report_json"

  # Append to CSV
  jq -r --arg r "$rate" '
    [ $r,
      .throughput,
      (.duration/1e9),
      .requests,
      (.status_codes["200"] // 0),
      (.requests - (.status_codes["200"] // 0)),
      ((1 - .success) * 100),
      (.latencies["50th"]/1e6),
      (.latencies["95th"]/1e6),
      (.latencies["99th"]/1e6),
      (.latencies.max/1e6)
    ] | @csv' "$report_json" >> "$RESULTS_CSV"
done

echo ""
echo "==============================="
echo " Run #${RUN_NUM} complete"
echo " Results: ${RESULTS_CSV}"
echo "==============================="

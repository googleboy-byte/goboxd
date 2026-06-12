import csv
import sys
import os
import matplotlib.pyplot as plt

# Accept run directory as argument, or default to latest run_N
base_dir = "docs/loadtest"

if len(sys.argv) > 1:
    run_dir = sys.argv[1]
else:
    # Find latest run_N directory
    run_dirs = sorted(
        [d for d in os.listdir(base_dir) if d.startswith("run_") and os.path.isdir(os.path.join(base_dir, d))],
        key=lambda d: int(d.split("_")[1])
    )
    if not run_dirs:
        print("No run directories found. Run load-test.sh first.")
        sys.exit(1)
    run_dir = os.path.join(base_dir, run_dirs[-1])

results_file = os.path.join(run_dir, "results.csv")
if not os.path.exists(results_file):
    print(f"No results.csv found in {run_dir}")
    sys.exit(1)

print(f"Plotting from: {results_file}")

rows = list(csv.DictReader(open(results_file)))
rps = [float(r["target_rps"]) for r in rows]

# Find breaking point (first RPS with error rate > 10%)
breaking_rps = None
for r in rows:
    if float(r["error_pct"]) > 10.0:
        breaking_rps = float(r["target_rps"])
        break

# Plot 1: Breaking Point
plt.figure()
plt.plot(rps, [float(r["error_pct"]) for r in rows], marker="o", color="red", label="Error Rate")
if breaking_rps is not None:
    plt.axvline(x=breaking_rps, color="blue", linestyle="--", label=f"Breaking Point ({int(breaking_rps)} RPS)")
plt.xlabel("Offered RPS")
plt.ylabel("Error rate (%)")
plt.title("Breaking point")
plt.legend()
plt.grid(True)
out1 = os.path.join(run_dir, "breaking-point.png")
plt.savefig(out1, dpi=150, bbox_inches="tight")
plt.close()

# Plot 2: Latency (p50, p95, p99)
plt.figure()
for k, lbl in [("p50_ms", "p50"), ("p95_ms", "p95"), ("p99_ms", "p99")]:
    plt.plot(rps, [float(r[k]) for r in rows], marker="o", label=lbl)
if breaking_rps is not None:
    plt.axvline(x=breaking_rps, color="blue", linestyle="--", label=f"Breaking Point ({int(breaking_rps)} RPS)")
plt.xlabel("Offered RPS")
plt.ylabel("Latency (ms)")
plt.title("RPS vs Latency")
plt.legend()
plt.grid(True)
out2 = os.path.join(run_dir, "latency.png")
plt.savefig(out2, dpi=150, bbox_inches="tight")
plt.close()

print(f"Plots saved to {out1} and {out2}")

# cgroup memory tracking

## how it works

The kernel maintains a per-cgroup memory counter in `mm/memcontrol.c`.
Every page fault that allocates a physical page increments the counter
for that cgroup and all ancestors. `memory.peak` is a high watermark —
it never decrements, making it the right file for peak RSS measurement.

## what nsjail creates

Each nsjail invocation creates a cgroup at `/sys/fs/cgroup/NSJAIL.<pid>/`
where `<pid>` is nsjail's pid as seen from the host cgroup namespace.
The sandboxed child process is moved into this cgroup, not nsjail itself.
This is why reading `cgroup.procs` and matching against `cmd.Process.Pid`
(the nsjail pid) never worked — the procs file contains the child's pid,
not nsjail's. (finding this out....took some time)

## why post-wait reads return 0

Nsjail cleans up its cgroup immediately on exit. By the time `cmd.Wait()`
returns in Go, the `NSJAIL.<pid>` directory is already gone. Any read
after wait gets file-not-found. 

note: nothing to read after it runs. had to track while running: 

```
docker exec goboxd sh -c 'ls /sys/fs/cgroup/ | grep NSJAIL'
docker exec goboxd sh -c 'cat /sys/fs/cgroup/NSJAIL.*/memory.peak'

NSJAIL.61
40894464
```

## why --cgroupns=host is required

Without `--cgroupns=host`, the container has its own cgroup namespace.
The memory controller charges land in the container's root cgroup
(`goboxd-node`), not in the per-nsjail child cgroups. The `NSJAIL.*`
directories exist but `memory.peak` reads 0 because memory accounting
is happening at a higher level in the hierarchy.

With `--cgroupns=host`, the container sees the real host cgroup hierarchy
and memory charges propagate correctly into `NSJAIL.<pid>/memory.peak`. (without this, could not access the memory.peak file)

## the polling approach

goboxd starts a goroutine before `cmd.Start()` that polls
`/sys/fs/cgroup/NSJAIL.*/memory.peak` every 10ms during execution.
It takes the max across all NSJAIL directories found (safe under the
semaphore-bounded concurrency limit). The goroutine is cancelled after
`cmd.Wait()` returns. Since `memory.peak` is a high watermark, even
intermittent reads capture the true peak.

## hybrid memory_exceeded detection

The system uses a two-pronged approach to detect memory exhaustion:

1.  **Threshold check**: If a process fails and `memory.peak` is >= 95% of the configured limit, it's marked as `memory_exceeded`. This works for processes that gradually consume memory.
2.  **Signature scanning**: Many languages (C++, Python, Java) throw exceptions upon allocation failure before the kernel formally charges the memory to the cgroup. The runner scans `stderr` for signatures like `bad_alloc`, `MemoryError`, or `OutOfMemoryError` to catch these cases early.

This hybrid model ensures robust OOM reporting even when the cgroup metrics aren't high enough to trigger a threshold alert.

## docker run requirement

The container must be started with `--cgroupns=host`:

    docker run -d --privileged --cgroupns=host --name goboxd -p 8080:8080 goboxd:latest

This is reflected in the Makefile `run` target.
# Plan Evolution

This file tracks the evolution of the `goboxd` system design and technical pivots.

## Direct Python execution to nsjail-wrapped execution

**What we thought we'd do:**
Invoke Python directly via `exec.Command` for the initial prototype.

**What we actually did:**
Replaced with nsjail invocation via `buildNsjailArgs()` in `sandbox.go`.

**Why it changed:**
The prototype needed to close the loop first. Nsjail was added once the request/response cycle was working end-to-end, reducing the debug surface during initial development.

## Single health endpoint to /healthz + /readyz + /info

**What we thought we'd do:**
Only implement `/healthz` for Stage 1.

**What we actually did:**
Added `/readyz` with language probing and `/info` with stats.

**Why it changed:**
The specification required all three for Stage 2. We built them during Stage 1 hardening to avoid accumulating debt, and because `/readyz` caching became a meaningful engineering problem.

## No concurrency limit to semaphore with queue timeout

**What we thought we'd do:**
Handle concurrency at the OS level via nsjail process limits.

**What we actually did:**
Added an in-process semaphore with a configurable size and queue timeout.

**Why it changed:**
The spec explicitly requires bounded concurrency with queuing. A queue timeout was added after identifying that requests could queue indefinitely if all slots stayed busy.

## Build flags only to build and run flags both validated

**What we thought we'd do:**
Validate `build.flags` against an allowlist only.

**What we actually did:**
Added `run.flag_allowlist` to the language config and validation path.

**Why it changed:**
The spec states that flags on both build and run must be filtered. Run flag injection is a real attack surface - a flag to the interpreter could redirect output or load external modules.

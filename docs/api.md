# API Documentation

The `goboxd` service exposes a REST API for executing untrusted code in secure sandboxes.

## Endpoints

### `POST /run`
Executes code across multiple test cases.

**Method:** `POST`  
**Path:** `/run`

#### Request Body
```json
{
  "language": "cpp",
  "source": "#include <iostream>\nint main(){std::cout<<\"hi\";}",
  "source_filename": "solution.cpp",
  "artifact_filename": "solution",
  "build": {
    "limits": { "wall_time_s": 5, "memory_kb": 1048576, "max_processes": 100 },
    "flags": ["-O2"]
  },
  "run": {
    "limits": { "wall_time_s": 3, "memory_kb": 524288, "max_processes": 64 },
    "flags": []
  },
  "tests": [
    { "stdin": "1\n", "expected_stdout": "hi" }
  ]
}
```

**Field Rules:**
- `language`: required. Must match a registered language ID.
- `source`: required. UTF-8 string, max 256 KiB.
- `source_filename`, `artifact_filename`: optional. Required for languages that use them (e.g., C++, Java). Must be a single path component, no separators, no leading dot, max 64 chars.
- `build`, `run`: optional. Override language defaults for limits and extra flags.
- `tests`: required. At least one test case.

#### Response Body (200 OK)
Returns 200 even if code fails to build or run.

```json
{
  "status": "wrong_output",
  "build": {
    "status": "ok",
    "stdout": "",
    "stderr": "",
    "duration_ms": 412
  },
  "tests": [
    {
      "status": "wrong_output",
      "stdout": "HI",
      "stderr": "",
      "duration_ms": 38,
      "memory_peak_kb": 8192
    }
  ]
}
```

#### Status Vocabulary
- **Top-level `status`**: `accepted` only if build is `ok` and every test is `accepted`. Otherwise, it is the first non-accepted status found in test order.
- **`build.status`**: `ok`, `failed`, `internal_error`.
- **`tests[].status`**: `accepted`, `wrong_output`, `output_whitespace_mismatch`, `time_exceeded`, `memory_exceeded`, `runtime_error`, `not_executed`, `internal_error`.

#### Error Responses (400 Bad Request)
Returned for malformed JSON, unknown languages, disallowed flags, or validation failures.
```json
{
  "error": {
    "code": "disallowed_flag",
    "message": "invalid flag: -fplugin=evil.so"
  }
}
```

---

### `GET /healthz`
Liveness check.
**Method:** `GET`  
**Path:** `/healthz`  
**Response:** `200 OK {"status":"ok"}`

---

### `GET /readyz`
Readiness check. Ensures that the sandboxing environment and all language toolchains are healthy.

**Method:** `GET`  
**Path:** `/readyz`

#### Response Body
- **Status 200 OK**: If nsjail and all languages are healthy.
- **Status 503 Service Unavailable**: If nsjail is broken or any language probe fails.

```json
{
  "status": "ok",
  "nsjail": {
    "ok": true,
    "version": "3.4"
  },
  "languages": {
    "bash": { "ok": true, "version": "GNU bash, version 5.2.15(1)-release (x86_64-pc-linux-gnu)" },
    "cpp": { "ok": true, "version": "g++ (Debian 12.2.0-14+deb12u1) 12.2.0" },
    "py3": { "ok": true, "version": "Python 3.11.2" },
    "rust": { "ok": true, "version": "rustc 1.63.0" }
  }
}
```

---

### `GET /info`
Service information and diagnostic metrics.

**Method:** `GET`  
**Path:** `/info`

#### Response Body (200 OK)
```json
{
  "build_info": {
    "version": "0.1.0",
    "commit": "4248131",
    "go_version": "go1.23.0"
  },
  "nsjail": {
    "path": "/usr/sbin/nsjail",
    "version": "3.4"
  },
  "languages": [
    {
      "id": "py3",
      "name": "Python 3",
      "version": "Python 3.11.2",
      "default_run_limits": { "wall_time_s": 9, "memory_kb": 102400, "max_processes": 100 }
    }
  ],
  "limits": {
    "max_source_bytes": 262144,
    "max_tests": 50,
    "max_concurrent_jobs": 4
  },
  "stats": {
    "in_flight_jobs": 0,
    "jobs_total": 150,
    "jobs_failed_internal": 2,
    "last_internal_error_at": "2026-05-21T16:50:00Z",
    "disk_free_bytes_jail_dir": 12685746176
  }
}
```
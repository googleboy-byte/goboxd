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

### `GET /readyz` (NOT YET IMPLEMENTED)
Readiness check.
**Intended Behavior:** Returns 200 only if `nsjail` is available and all language toolchains pass a smoke test. Returns 503 otherwise.

---

### `GET /info` (NOT YET IMPLEMENTED)
Service information.
**Intended Behavior:** Returns versions, registered languages, configured limits, and runtime statistics.

# Test Coverage Index

This document maps system infrastructure and logic to its corresponding test coverage.

## Unit Tests (`tests/unit/`)

| Component | Target File | Test File | Covered Logic |
|---|---|---|---|
| **Configuration** | [config.go](file:///home/violet/Desktop/goboxd/internal/config/config.go) | [config_test.go](file:///home/violet/Desktop/goboxd/tests/unit/config_test.go) | YAML parsing, empty list validation, missing ID detection, language lookup. |
| **Validation** | [validate.go](file:///home/violet/Desktop/goboxd/internal/validate/validate.go) | [validate_test.go](file:///home/violet/Desktop/goboxd/tests/unit/validate_test.go) | Filename security, glob flag allowlisting, source size enforcement, test count limits. |
| **Request Handling** | [run.go](file:///home/violet/Desktop/goboxd/internal/handler/run.go) | [handler_test.go](file:///home/violet/Desktop/goboxd/tests/unit/handler_test.go) | HTTP 400 paths: `invalid_json`, `unknown_language`, `disallowed_flag`, `bad_request` (size/tests). |
| **Utilities** | [runner.go](file:///home/violet/Desktop/goboxd/internal/runner/runner.go) | [resolve_test.go](file:///home/violet/Desktop/goboxd/tests/unit/resolve_test.go) | `ResolveString` and `ResolveArgs` placeholder substitution. |

## Integration Tests (`tests/integration/`)

| Test Area | Source | Target | Covered Logic |
|---|---|---|---|
| **Language Smoke Test** | [run_all.sh](file:///home/violet/Desktop/goboxd/tests/integration/run_all.sh) | API + Nsjail | End-to-end execution of Python, C++, Bash, and Rust solutions. |
| **Sandbox Execution** | [runner_integration_test.go](file:///home/violet/Desktop/goboxd/tests/integration/runner_integration_test.go) | [runner.go](file:///home/violet/Desktop/goboxd/internal/runner/runner.go) | `runner.Run` logic, literal matches, and whitespace sensitivity (requires `nsjail`). |

## Key Coverage Paths

### 1. Security & Validation
- **Path**: `POST /run` → `handler.NewRunHandler` → `validate.ValidateRunRequest`
- **Tests**: `TestRunHandler_OversizeBody`, `TestRunHandler_TooManyTests`, `TestValidateRunRequest/Too_much_source`.

### 2. Flag Hardening
- **Path**: `POST /run` → `handler.NewRunHandler` → `validate.ValidateFlags`
- **Tests**: `TestRunHandler_DisallowedFlag`, `TestValidateFlags/Not_allowed`, `TestValidateFlags/Empty_allowlist`.

### 3. Service Diagnostics
- **Path**: `GET /info` / `GET /readyz` → `handler.HealthHandler`
- **Verification**: Manual via `curl` as specified in Stage 1 requirements.

# goboxd

This is a sandbox daemon written in go. It is used to run untrusted code in a sandbox.

## Run

```
make run            # builds image, starts container on port 8080
make test           # runs unit tests
make integration    # end-to-end tests, requires running container
make lint           # runs static analysis
```

## Docs

```
docs/api.md - HTTP contract
docs/languages.md - supported languages and YAML schema
docs/testing.md - test coverage index
docs/architecture.md - system architecture
docs/benchmarks.md - benchmarks
docs/security.md - security considerations
docs/test_logs.md - test logs
```

## Framework

HTTP routing uses chi, stays close to net/http, making handlers trivial to test with httptest and avoiding framework lock-in.
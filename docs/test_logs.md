# End-to-End Test Logs

This file maintains a record of successful and failed execution runs for all supported languages.

## C++ (cpp)

### [SUCCESS] Standard Hello World
**Request:**
```bash
curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"cpp",
  "source":"#include <iostream>\nint main(){std::cout<<\"hi\";return 0;}",
  "tests":[{"stdin":"","expected_stdout":"hi"}]
}' http://localhost:8080/run
```
**Response:**
```json
{"status":"accepted","build":{"status":"ok","stdout":"","stderr":"","duration_ms":329},"tests":[{"status":"accepted","stdout":"hi","stderr":"","duration_ms":5,"memory_peak_kb":0}]}
```

### [FAILURE] Syntax Error
**Request:**
```bash
curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"cpp",
  "source":"this is not valid c++",
  "tests":[{"stdin":"","expected_stdout":"hi"}]
}' http://localhost:8080/run
```
**Response:**
```json
{"status":"build_failed","build":{"status":"failed","stdout":"","stderr":"/sandbox/solution.cpp:1:1: error: expected unqualified-id before 'this'\n    1 | this is not valid c++\n      | ^~~~\n","duration_ms":14},"tests":[{"status":"not_executed","stdout":"","stderr":"","duration_ms":0,"memory_peak_kb":0}]}
```

## Python 3 (py3)

### [SUCCESS] Standard Hello World (with Whitespace Mismatch)
**Request:**
```bash
curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"py3",
  "source":"print(\"hello\")",
  "tests":[{"stdin":"","expected_stdout":"hello"}]
}' http://localhost:8080/run
```
**Response:**
```json
{"status":"output_whitespace_mismatch","tests":[{"status":"output_whitespace_mismatch","stdout":"hello\n","stderr":"","duration_ms":34,"memory_peak_kb":0}]}
```

### [FAILURE] Runtime Error
**Request:**
```bash
curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"py3",
  "source":"1/0",
  "tests":[{"stdin":"","expected_stdout":"hello"}]
}' http://localhost:8080/run
```
**Response:**
```json
{"status":"runtime_error","tests":[{"status":"runtime_error","stdout":"","stderr":"Traceback (most recent call last):\n  File \"/sandbox/solution.py\", line 1, in <module>\n    1/0\n    ~^~\nZeroDivisionError: division by zero\n","duration_ms":15,"memory_peak_kb":0}]}
```

## Bash (bash)

### [SUCCESS] echo hello
**Request:**
```bash
curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"bash",
  "source":"echo hello",
  "tests":[{"stdin":"","expected_stdout":"hello\n"}]
}' http://localhost:8080/run
```
**Response:**
```json
{"status":"accepted","tests":[{"status":"accepted","stdout":"hello\n","stderr":"","duration_ms":7,"memory_peak_kb":0}]}
```

### [FAILURE] exit 1
**Request:**
```bash
curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"bash",
  "source":"exit 1",
  "tests":[{"stdin":"","expected_stdout":""}]
}' http://localhost:8080/run
```
**Response:**
```json
{"status":"runtime_error","tests":[{"status":"runtime_error","stdout":"","stderr":"","duration_ms":5,"memory_peak_kb":0}]}
```

## Rust (rust)

### [SUCCESS] Standard Hello World
**Request:**
```bash
curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"rust",
  "source":"fn main() { println!(\"hi\"); }",
  "tests":[{"stdin":"","expected_stdout":"hi\n"}]
}' http://localhost:8080/run
```
**Response:**
```json
{"status":"accepted","build":{"status":"ok","stdout":"","stderr":"","duration_ms":285},"tests":[{"status":"accepted","stdout":"hi\n","stderr":"","duration_ms":7,"memory_peak_kb":0}]}
```

### [FAILURE] Syntax Error
**Request:**
```bash
curl -s -X POST -H "Content-Type: application/json" -d '{
  "language":"rust",
  "source":"fn main() { invalid code }",
  "tests":[{"stdin":"","expected_stdout":""}]
}' http://localhost:8080/run
```
**Response:**
```json
{"status":"build_failed","build":{"status":"failed","stdout":"","stderr":"error: expected one of `!`, `.`, `::`, `;`, `?`, `{`, `}`, or an operator, found `code`...","duration_ms":53},"tests":[{"status":"not_executed","stdout":"","stderr":"","duration_ms":0,"memory_peak_kb":0}]}
```
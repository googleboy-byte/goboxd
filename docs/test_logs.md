# Final Project Verification Log (16 Languages) - Thu May 28 15:47:24 IST 2026

## 1. Make Lint
Linting code...
go vet ./...
staticcheck not found, skipping (go vet passed)
## 2. Make Test
Running unit tests...
go test -v ./tests/unit/...
=== RUN   TestConfigLoad
=== RUN   TestConfigLoad/Valid_file_loads_py3_correctly
=== RUN   TestConfigLoad/Missing_file_returns_an_error
=== RUN   TestConfigLoad/Unknown_language_returns_an_error
=== RUN   TestConfigLoad/Bad_YAML_fails
=== RUN   TestConfigLoad/Empty_languages_list_returns_an_error
=== RUN   TestConfigLoad/Language_missing_ID_returns_an_error
=== RUN   TestConfigLoad/Language_with_zero_wall_time_s_returns_an_error
=== RUN   TestConfigLoad/Language_with_missing_limits_returns_an_error
--- PASS: TestConfigLoad (0.00s)
    --- PASS: TestConfigLoad/Valid_file_loads_py3_correctly (0.00s)
    --- PASS: TestConfigLoad/Missing_file_returns_an_error (0.00s)
    --- PASS: TestConfigLoad/Unknown_language_returns_an_error (0.00s)
    --- PASS: TestConfigLoad/Bad_YAML_fails (0.00s)
    --- PASS: TestConfigLoad/Empty_languages_list_returns_an_error (0.00s)
    --- PASS: TestConfigLoad/Language_missing_ID_returns_an_error (0.00s)
    --- PASS: TestConfigLoad/Language_with_zero_wall_time_s_returns_an_error (0.00s)
    --- PASS: TestConfigLoad/Language_with_missing_limits_returns_an_error (0.00s)
=== RUN   TestRunHandler_UnknownLanguage
2026/05/28 15:43:40 INFO request completed request_id=a2400f4f0f1f9b2d language=fortran status=unknown_language duration_ms=0
--- PASS: TestRunHandler_UnknownLanguage (0.00s)
=== RUN   TestRunHandler_MalformedJSON
2026/05/28 15:43:40 INFO request completed request_id=ae9c9a7e5df66cae language=unknown status=invalid_json duration_ms=0
--- PASS: TestRunHandler_MalformedJSON (0.00s)
=== RUN   TestRunHandler_DisallowedFlag
2026/05/28 15:43:40 INFO request completed request_id=cac144b284a2f844 language=cpp status=disallowed_flag duration_ms=0
--- PASS: TestRunHandler_DisallowedFlag (0.00s)
=== RUN   TestRunHandler_EmptySource
2026/05/28 15:43:40 INFO request completed request_id=3c5317d9eda2ca60 language=py3 status=bad_request duration_ms=0
--- PASS: TestRunHandler_EmptySource (0.00s)
=== RUN   TestRunHandler_TooManyTests
2026/05/28 15:43:40 INFO request completed request_id=a05db40664ae9a8c language=py3 status=bad_request duration_ms=0
--- PASS: TestRunHandler_TooManyTests (0.00s)
=== RUN   TestRunHandler_NoTests
2026/05/28 15:43:40 INFO request completed request_id=d7aae57dfe2b5ea9 language=py3 status=bad_request duration_ms=0
--- PASS: TestRunHandler_NoTests (0.00s)
=== RUN   TestRunHandler_OversizeBody
2026/05/28 15:43:40 INFO request completed request_id=87bcea81b8d2b4c1 language=py3 status=bad_request duration_ms=1
--- PASS: TestRunHandler_OversizeBody (0.00s)
=== RUN   TestRunHandler_QueueTimeout
2026/05/28 15:43:41 INFO request completed request_id=4730155af56a1a91 language=unknown status=queue_timeout duration_ms=1000
--- PASS: TestRunHandler_QueueTimeout (1.00s)
=== RUN   TestResolveString
=== RUN   TestResolveString/Single_placeholder
=== RUN   TestResolveString/Multiple_placeholders
=== RUN   TestResolveString/No_placeholders
=== RUN   TestResolveString/Unknown_placeholder_left_as-is
=== RUN   TestResolveString/Empty_string
--- PASS: TestResolveString (0.00s)
    --- PASS: TestResolveString/Single_placeholder (0.00s)
    --- PASS: TestResolveString/Multiple_placeholders (0.00s)
    --- PASS: TestResolveString/No_placeholders (0.00s)
    --- PASS: TestResolveString/Unknown_placeholder_left_as-is (0.00s)
    --- PASS: TestResolveString/Empty_string (0.00s)
=== RUN   TestResolveArgs
=== RUN   TestResolveArgs/Resolves_all_args
=== RUN   TestResolveArgs/Empty_args
=== RUN   TestResolveArgs/No_placeholders_in_args
--- PASS: TestResolveArgs (0.00s)
    --- PASS: TestResolveArgs/Resolves_all_args (0.00s)
    --- PASS: TestResolveArgs/Empty_args (0.00s)
    --- PASS: TestResolveArgs/No_placeholders_in_args (0.00s)
=== RUN   TestValidateFilename
=== RUN   TestValidateFilename/Happy_path
=== RUN   TestValidateFilename/Single_component
=== RUN   TestValidateFilename/Empty
=== RUN   TestValidateFilename/Too_long
=== RUN   TestValidateFilename/Path_traversal_..
=== RUN   TestValidateFilename/Path_separator_/
=== RUN   TestValidateFilename/Reserved_.
=== RUN   TestValidateFilename/Reserved_..
=== RUN   TestValidateFilename/Hidden_file
--- PASS: TestValidateFilename (0.00s)
    --- PASS: TestValidateFilename/Happy_path (0.00s)
    --- PASS: TestValidateFilename/Single_component (0.00s)
    --- PASS: TestValidateFilename/Empty (0.00s)
    --- PASS: TestValidateFilename/Too_long (0.00s)
    --- PASS: TestValidateFilename/Path_traversal_.. (0.00s)
    --- PASS: TestValidateFilename/Path_separator_/ (0.00s)
    --- PASS: TestValidateFilename/Reserved_. (0.00s)
    --- PASS: TestValidateFilename/Reserved_.. (0.00s)
    --- PASS: TestValidateFilename/Hidden_file (0.00s)
=== RUN   TestValidateFlags
=== RUN   TestValidateFlags/Empty_requested
=== RUN   TestValidateFlags/All_allowed
=== RUN   TestValidateFlags/Not_allowed
=== RUN   TestValidateFlags/Mixed
=== RUN   TestValidateFlags/Empty_allowlist
--- PASS: TestValidateFlags (0.00s)
    --- PASS: TestValidateFlags/Empty_requested (0.00s)
    --- PASS: TestValidateFlags/All_allowed (0.00s)
    --- PASS: TestValidateFlags/Not_allowed (0.00s)
    --- PASS: TestValidateFlags/Mixed (0.00s)
    --- PASS: TestValidateFlags/Empty_allowlist (0.00s)
=== RUN   TestValidateRunRequest
=== RUN   TestValidateRunRequest/Happy_path
=== RUN   TestValidateRunRequest/Empty_language
=== RUN   TestValidateRunRequest/Empty_source
=== RUN   TestValidateRunRequest/Too_much_source
=== RUN   TestValidateRunRequest/Too_many_tests
=== RUN   TestValidateRunRequest/Zero_tests
--- PASS: TestValidateRunRequest (0.00s)
    --- PASS: TestValidateRunRequest/Happy_path (0.00s)
    --- PASS: TestValidateRunRequest/Empty_language (0.00s)
    --- PASS: TestValidateRunRequest/Empty_source (0.00s)
    --- PASS: TestValidateRunRequest/Too_much_source (0.00s)
    --- PASS: TestValidateRunRequest/Too_many_tests (0.00s)
    --- PASS: TestValidateRunRequest/Zero_tests (0.00s)
PASS
ok  	github.com/thesouldev/goboxd/tests/unit	(cached)
## 3. Make Secure
Running security verification tests...
Waiting for http://localhost:8080 to be ready...
Server is ready!
Verifying security holes for http://localhost:8080...
[Hole 1] Path Traversal via SourceFilename...
  PASS: Rejected malicious SourceFilename
[Hole 1] Path Traversal via ArtifactFilename...
  PASS: Rejected malicious ArtifactFilename
[Hole 3] Compiler Flag Injection...
  PASS: Rejected disallowed build flag
[Hole 3] Run Flag Injection...
  PASS: Rejected disallowed run flag
[Hole 4] Request Size Limits...
  PASS: Rejected large source
  PASS: Rejected large stdin
[Hole 6] Output Truncation Marker...
  PASS: Truncation marker present
[Hole 7] Network Isolation...
  PASS: Network is isolated
ALL SECURITY TESTS PASSED
## 4. Make Integration (All 16)
Running integration tests...
bash tests/integration/run_all.sh http://localhost:8080
Discovering registered languages...
--- Integration Tests ---
Testing py3...
✅ py3: accepted
Testing cpp...
✅ cpp: accepted
Testing bash...
✅ bash: accepted
Testing rust...
✅ rust: accepted
Testing java...
✅ java: accepted
Testing c...
✅ c: accepted
Testing js...
✅ js: accepted
Testing verilog...
✅ verilog: accepted
Testing go...
✅ go: accepted
Testing kotlin...
✅ kotlin: accepted
Testing csharp...
✅ csharp: accepted
Testing ruby...
✅ ruby: accepted
Testing lua...
✅ lua: accepted
Testing ocaml...
✅ ocaml: accepted
Testing swift...
✅ swift: accepted
Testing zig...
✅ zig: accepted
--- Integration tests completed! ---
## 5. Make Corpus (Full Run)
bash tests/corpus/run_corpus.sh http://localhost:8080
==============================
 goboxd corpus test suite
 Server: http://localhost:8080
==============================

── Waiting for server ──
[32m✅ server reachable[0m

── Health endpoints ──
[32m✅ /healthz returns ok[0m
[32m✅ /readyz returns ok[0m
[32m✅ /info has build_info.version[0m
[32m✅ /info has languages[0m
[32m✅ /info has limits.max_source_bytes[0m
[32m✅ /info has stats.jobs_total[0m

── Happy path: hello world ──
[32m✅ py3 hello world[0m
[32m✅ cpp hello world[0m
[32m✅ c hello world[0m
[32m✅ bash hello world[0m
[32m✅ js hello world[0m
[32m✅ rust hello world[0m
[32m✅ java hello world[0m
[32m✅ verilog hello world[0m

── stdin echo ──
[32m✅ py3 stdin echo[0m
[32m✅ c stdin echo[0m

── multiple test cases ──
[32m✅ py3 multi-test accepted[0m
[32m✅ py3 multi-test: 3 results returned[0m

── status: wrong_output ──
[32m✅ py3 wrong_output[0m
[32m✅ py3 wrong_output: build.status ok[0m

── status: output_whitespace_mismatch ──
[32m✅ py3 whitespace mismatch detected (output_whitespace_mismatch)[0m

── build failure → not_executed ──
[32m✅ cpp build_failed top-level[0m
[32m✅ cpp build.status failed[0m
[32m✅ cpp all tests not_executed after build failure[0m
[32m✅ py3 syntax error detected (runtime_error)[0m

── flag override ──
[32m✅ cpp with allowed flags[0m

── status: runtime_error ──
[32m✅ c abort() → runtime_error[0m
[32m✅ py3 exception → runtime_error[0m

── status: time_exceeded ──
[32m✅ py3 infinite loop → time_exceeded[0m
[32m✅ c infinite loop → time_exceeded[0m

── edge: empty output ──
[32m✅ py3 empty output accepted[0m

── edge: 50 test cases ──
[32m✅ py3 50 tests accepted[0m
[32m✅ py3 50 tests: 50 results returned[0m

── edge: source size boundary ──
[32m✅ py3 source at 256KiB-1 accepted[0m

── adversarial: path traversal ──
[32m✅ path traversal source_filename rejected[0m
[32m✅ path traversal artifact_filename rejected[0m

── adversarial: disallowed flags ──
[32m✅ cpp disallowed build flag rejected[0m
[32m✅ cpp --specs flag rejected[0m

── adversarial: oversize body ──
[32m✅ oversize body rejected[0m

── adversarial: unknown language ──
[32m✅ unknown language rejected[0m

── adversarial: malformed JSON ──
[32m✅ malformed JSON rejected[0m

── adversarial: missing fields ──
[32m✅ missing source rejected[0m
[32m✅ missing language rejected[0m
[32m✅ missing tests rejected[0m
[32m✅ empty tests array rejected[0m

── adversarial: Java without required filenames ──
[32m✅ java without source_filename rejected[0m
[32m✅ java without artifact_filename rejected[0m

── adversarial: too many tests ──
[32m✅ 51 tests rejected[0m

── load: 200 requests at c=50 ──
[32m✅ load: all 200 requests returned 200[0m

==============================
 Results: 50/50 passed (0 skipped)
[32m✅  ALL CORPUS TESTS PASSED[0m
## 6. Make Load
bash tests/load/load.sh http://localhost:8080
--- Concurrency 1 ---

Summary:
  Total:	3.5997 secs
  Slowest:	0.0459 secs
  Fastest:	0.0151 secs
  Average:	0.0180 secs
  Requests/sec:	55.5603
  
  Total data:	36800 bytes
  Size/request:	184 bytes

Response time histogram:
  0.015 [1]	|
  0.018 [168]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.021 [17]	|■■■■
  0.024 [4]	|■
  0.027 [2]	|
  0.031 [3]	|■
  0.034 [1]	|
  0.037 [0]	|
  0.040 [2]	|
  0.043 [0]	|
  0.046 [2]	|


Latency distribution:
  10%% in 0.0159 secs
  25%% in 0.0167 secs
  50%% in 0.0171 secs
  75%% in 0.0176 secs
  90%% in 0.0192 secs
  95%% in 0.0252 secs
  99%% in 0.0448 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0000 secs, 0.0004 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0002 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0001 secs
  resp wait:	0.0179 secs, 0.0150 secs, 0.0458 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0003 secs

Status code distribution:
  [200]	200 responses



--- Concurrency 10 ---

Summary:
  Total:	1.8085 secs
  Slowest:	0.2024 secs
  Fastest:	0.0216 secs
  Average:	0.0869 secs
  Requests/sec:	110.5917
  
  Total data:	36801 bytes
  Size/request:	184 bytes

Response time histogram:
  0.022 [1]	|■
  0.040 [6]	|■■■
  0.058 [3]	|■■
  0.076 [48]	|■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.094 [75]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.112 [43]	|■■■■■■■■■■■■■■■■■■■■■■■
  0.130 [21]	|■■■■■■■■■■■
  0.148 [2]	|■
  0.166 [0]	|
  0.184 [0]	|
  0.202 [1]	|■


Latency distribution:
  10%% in 0.0646 secs
  25%% in 0.0747 secs
  50%% in 0.0842 secs
  75%% in 0.1010 secs
  90%% in 0.1159 secs
  95%% in 0.1224 secs
  99%% in 0.1370 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0000 secs, 0.0009 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0005 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0003 secs
  resp wait:	0.0867 secs, 0.0209 secs, 0.2023 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0021 secs

Status code distribution:
  [200]	200 responses



--- Concurrency 50 ---

Summary:
  Total:	2.0756 secs
  Slowest:	0.8649 secs
  Fastest:	0.0372 secs
  Average:	0.4691 secs
  Requests/sec:	96.3557
  
  Total data:	36805 bytes
  Size/request:	184 bytes

Response time histogram:
  0.037 [1]	|■
  0.120 [10]	|■■■■■■
  0.203 [10]	|■■■■■■
  0.285 [9]	|■■■■■■
  0.368 [30]	|■■■■■■■■■■■■■■■■■■
  0.451 [65]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.534 [20]	|■■■■■■■■■■■■
  0.617 [5]	|■■■
  0.699 [2]	|■
  0.782 [14]	|■■■■■■■■■
  0.865 [34]	|■■■■■■■■■■■■■■■■■■■■■


Latency distribution:
  10%% in 0.2012 secs
  25%% in 0.3553 secs
  50%% in 0.4129 secs
  75%% in 0.6623 secs
  90%% in 0.8289 secs
  95%% in 0.8448 secs
  99%% in 0.8648 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0003 secs, 0.0000 secs, 0.0046 secs
  DNS-lookup:	0.0002 secs, 0.0000 secs, 0.0020 secs
  req write:	0.0001 secs, 0.0000 secs, 0.0019 secs
  resp wait:	0.4686 secs, 0.0368 secs, 0.8647 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0032 secs

Status code distribution:
  [200]	200 responses



--- Concurrency 100 ---

Summary:
  Total:	1.7536 secs
  Slowest:	0.9927 secs
  Fastest:	0.0416 secs
  Average:	0.6912 secs
  Requests/sec:	114.0518
  
  Total data:	36800 bytes
  Size/request:	184 bytes

Response time histogram:
  0.042 [1]	|■
  0.137 [13]	|■■■■■■■■■
  0.232 [10]	|■■■■■■■
  0.327 [10]	|■■■■■■■
  0.422 [11]	|■■■■■■■
  0.517 [7]	|■■■■■
  0.612 [9]	|■■■■■■
  0.707 [13]	|■■■■■■■■■
  0.803 [11]	|■■■■■■■
  0.898 [59]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.993 [56]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■


Latency distribution:
  10%% in 0.1938 secs
  25%% in 0.5085 secs
  50%% in 0.8331 secs
  75%% in 0.9123 secs
  90%% in 0.9422 secs
  95%% in 0.9509 secs
  99%% in 0.9770 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0040 secs, 0.0000 secs, 0.0184 secs
  DNS-lookup:	0.0012 secs, 0.0000 secs, 0.0041 secs
  req write:	0.0017 secs, 0.0000 secs, 0.0081 secs
  resp wait:	0.6853 secs, 0.0345 secs, 0.9924 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0046 secs

Status code distribution:
  [200]	200 responses




# Final Project Verification Log - Wed May 27 23:03:43 IST 2026

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
2026/05/26 23:02:00 INFO request completed request_id=cc97faf9c77e8c1e language=fortran status=unknown_language duration_ms=0
--- PASS: TestRunHandler_UnknownLanguage (0.00s)
=== RUN   TestRunHandler_MalformedJSON
2026/05/26 23:02:00 INFO request completed request_id=9541d3a7cb27778c language=unknown status=invalid_json duration_ms=0
--- PASS: TestRunHandler_MalformedJSON (0.00s)
=== RUN   TestRunHandler_DisallowedFlag
2026/05/26 23:02:00 INFO request completed request_id=c4cb3d8efe6c5e0e language=cpp status=disallowed_flag duration_ms=0
--- PASS: TestRunHandler_DisallowedFlag (0.00s)
=== RUN   TestRunHandler_EmptySource
2026/05/26 23:02:00 INFO request completed request_id=d2d754ef100e07bd language=py3 status=bad_request duration_ms=0
--- PASS: TestRunHandler_EmptySource (0.00s)
=== RUN   TestRunHandler_TooManyTests
2026/05/26 23:02:00 INFO request completed request_id=71f273d7f6acfa77 language=py3 status=bad_request duration_ms=0
--- PASS: TestRunHandler_TooManyTests (0.00s)
=== RUN   TestRunHandler_NoTests
2026/05/26 23:02:00 INFO request completed request_id=8657337bcca8b73e language=py3 status=bad_request duration_ms=0
--- PASS: TestRunHandler_NoTests (0.00s)
=== RUN   TestRunHandler_OversizeBody
2026/05/26 23:02:00 INFO request completed request_id=124337019bcfa8b6 language=py3 status=bad_request duration_ms=1
--- PASS: TestRunHandler_OversizeBody (0.00s)
=== RUN   TestRunHandler_QueueTimeout
2026/05/26 23:02:01 INFO request completed request_id=333b5579a0ed2de8 language=unknown status=queue_timeout duration_ms=1000
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
make: *** [Makefile:15: secure] Interrupt
## 6. Make Load
bash tests/load/load.sh http://localhost:8080
--- Concurrency 1 ---

Summary:
  Total:	5.2939 secs
  Slowest:	1.7560 secs
  Fastest:	0.0147 secs
  Average:	0.0265 secs
  Requests/sec:	37.7793
  
  Total data:	36802 bytes
  Size/request:	184 bytes

Response time histogram:
  0.015 [1]	|
  0.189 [198]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.363 [0]	|
  0.537 [0]	|
  0.711 [0]	|
  0.885 [0]	|
  1.059 [0]	|
  1.234 [0]	|
  1.408 [0]	|
  1.582 [0]	|
  1.756 [1]	|


Latency distribution:
  10%% in 0.0150 secs
  25%% in 0.0158 secs
  50%% in 0.0164 secs
  75%% in 0.0173 secs
  90%% in 0.0197 secs
  95%% in 0.0286 secs
  99%% in 0.0588 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0000 secs, 0.0006 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0002 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0001 secs
  resp wait:	0.0264 secs, 0.0146 secs, 1.7550 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0007 secs

Status code distribution:
  [200]	200 responses



--- Concurrency 10 ---

Summary:
  Total:	1.6030 secs
  Slowest:	0.1739 secs
  Fastest:	0.0221 secs
  Average:	0.0787 secs
  Requests/sec:	124.7629
  
  Total data:	36801 bytes
  Size/request:	184 bytes

Response time histogram:
  0.022 [1]	|■
  0.037 [0]	|
  0.052 [5]	|■■■
  0.068 [53]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.083 [78]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.098 [37]	|■■■■■■■■■■■■■■■■■■■
  0.113 [12]	|■■■■■■
  0.128 [7]	|■■■■
  0.144 [3]	|■■
  0.159 [3]	|■■
  0.174 [1]	|■


Latency distribution:
  10%% in 0.0587 secs
  25%% in 0.0655 secs
  50%% in 0.0742 secs
  75%% in 0.0873 secs
  90%% in 0.1035 secs
  95%% in 0.1197 secs
  99%% in 0.1581 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0001 secs, 0.0000 secs, 0.0022 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0008 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0011 secs
  resp wait:	0.0785 secs, 0.0212 secs, 0.1738 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0009 secs

Status code distribution:
  [200]	200 responses



--- Concurrency 50 ---

Summary:
  Total:	2.0587 secs
  Slowest:	0.6923 secs
  Fastest:	0.0372 secs
  Average:	0.4461 secs
  Requests/sec:	97.1464
  
  Total data:	36805 bytes
  Size/request:	184 bytes

Response time histogram:
  0.037 [1]	|■
  0.103 [11]	|■■■■■■
  0.168 [9]	|■■■■■
  0.234 [8]	|■■■■
  0.299 [11]	|■■■■■■
  0.365 [11]	|■■■■■■
  0.430 [13]	|■■■■■■■
  0.496 [19]	|■■■■■■■■■■■
  0.561 [72]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.627 [40]	|■■■■■■■■■■■■■■■■■■■■■■
  0.692 [5]	|■■■


Latency distribution:
  10%% in 0.1643 secs
  25%% in 0.3639 secs
  50%% in 0.5239 secs
  75%% in 0.5581 secs
  90%% in 0.5787 secs
  95%% in 0.5901 secs
  99%% in 0.6575 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0005 secs, 0.0000 secs, 0.0040 secs
  DNS-lookup:	0.0003 secs, 0.0000 secs, 0.0023 secs
  req write:	0.0002 secs, 0.0000 secs, 0.0012 secs
  resp wait:	0.4454 secs, 0.0362 secs, 0.6922 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0016 secs

Status code distribution:
  [200]	200 responses



--- Concurrency 100 ---

Summary:
  Total:	2.5455 secs
  Slowest:	1.5187 secs
  Fastest:	0.0408 secs
  Average:	0.9008 secs
  Requests/sec:	78.5713
  
  Total data:	36811 bytes
  Size/request:	184 bytes

Response time histogram:
  0.041 [1]	|■
  0.189 [15]	|■■■■■■■■■
  0.336 [15]	|■■■■■■■■■
  0.484 [14]	|■■■■■■■■
  0.632 [12]	|■■■■■■■
  0.780 [12]	|■■■■■■■
  0.928 [13]	|■■■■■■■■
  1.075 [12]	|■■■■■■■
  1.223 [66]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.371 [6]	|■■■■
  1.519 [34]	|■■■■■■■■■■■■■■■■■■■■■


Latency distribution:
  10%% in 0.2573 secs
  25%% in 0.5450 secs
  50%% in 1.0915 secs
  75%% in 1.1501 secs
  90%% in 1.4307 secs
  95%% in 1.4448 secs
  99%% in 1.5023 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0057 secs, 0.0000 secs, 0.0217 secs
  DNS-lookup:	0.0049 secs, 0.0000 secs, 0.0194 secs
  req write:	0.0001 secs, 0.0000 secs, 0.0011 secs
  resp wait:	0.8948 secs, 0.0401 secs, 1.5186 secs
  resp read:	0.0002 secs, 0.0000 secs, 0.0085 secs

Status code distribution:
  [200]	200 responses



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

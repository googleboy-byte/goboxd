mainak@mainak-HP-EliteBook-830-G5:~/Desktop/goboxd$ make payloads
bash tests/corpus/run_payloads.sh http://localhost:8080
==============================
 Running Custom Payloads
 Server: http://localhost:8080
==============================
Testing dart (accepted)... ✅ PASSED
Testing dart (runtime_error)... ✅ PASSED
Testing dart (time_exceeded)... ✅ PASSED
Testing dart (wrong_output)... ✅ PASSED
Testing fortran (accepted)... ✅ PASSED
Testing fortran (build_failed)... ✅ PASSED
Testing fortran (runtime_error)... ✅ PASSED
Testing fortran (time_exceeded)... ✅ PASSED
Testing fortran (wrong_output)... ✅ PASSED
Testing typescript (accepted)... ✅ PASSED
Testing typescript (build_failed)... ✅ PASSED
Testing typescript (runtime_error)... ✅ PASSED
Testing typescript (time_exceeded)... ✅ PASSED
Testing typescript (wrong_output)... ✅ PASSED

==============================
 Results: 14/14 passed
 ALL PAYLOADS PASSED
mainak@mainak-HP-EliteBook-830-G5:~/Desktop/goboxd$ make payloads
bash tests/corpus/run_payloads.sh http://localhost:8080
==============================
 Running Custom Payloads
 Server: http://localhost:8080
==============================
Testing dart (accepted)... ✅ PASSED
  Expected: accepted
  Got:      accepted
  Response: {"status":"accepted","build":{"status":"ok","stdout":"Generated: /sandbox/solution\n","stderr":"","duration_ms":2458},"tests":[{"status":"accepted","stdout":"42\n","stderr":"","duration_ms":31,"memory_peak_kb":812}]}
Testing dart (runtime_error)... ✅ PASSED
  Expected: runtime_error
  Got:      runtime_error
  Response: {"status":"runtime_error","build":{"status":"ok","stdout":"Generated: /sandbox/solution\n","stderr":"","duration_ms":2188},"tests":[{"status":"runtime_error","stdout":"","stderr":"Unhandled exception:\nboom\n#0      main (file:///sandbox/solution.dart:1)\n#1      _delayEntrypointInvocation.\u003canonymous closure\u003e (dart:isolate-patch/isolate_patch.dart:297)\n#2      _RawReceivePort._handleMessage (dart:isolate-patch/isolate_patch.dart:184)\n","duration_ms":16,"memory_peak_kb":512}]}
Testing dart (time_exceeded)... ✅ PASSED
  Expected: time_exceeded
  Got:      time_exceeded
  Response: {"status":"time_exceeded","build":{"status":"ok","stdout":"Generated: /sandbox/solution\n","stderr":"","duration_ms":2249},"tests":[{"status":"time_exceeded","stdout":"","stderr":"","duration_ms":1004,"memory_peak_kb":1836}]}
Testing dart (wrong_output)... ✅ PASSED
  Expected: wrong_output
  Got:      wrong_output
  Response: {"status":"wrong_output","build":{"status":"ok","stdout":"Generated: /sandbox/solution\n","stderr":"","duration_ms":2413},"tests":[{"status":"wrong_output","stdout":"42\n","stderr":"","duration_ms":18,"memory_peak_kb":256}]}
Testing fortran (accepted)... ✅ PASSED
  Expected: accepted
  Got:      accepted
  Response: {"status":"accepted","build":{"status":"ok","stdout":"","stderr":"","duration_ms":179},"tests":[{"status":"accepted","stdout":"42\n","stderr":"","duration_ms":18,"memory_peak_kb":256}]}
Testing fortran (build_failed)... ✅ PASSED
  Expected: build_failed
  Got:      build_failed
  Response: {"status":"build_failed","build":{"status":"failed","stdout":"","stderr":"f951: Error: Unexpected end of file in '/sandbox/solution.f90'\n","duration_ms":35},"tests":[{"status":"not_executed","stdout":"","stderr":"","duration_ms":0,"memory_peak_kb":0}]}
Testing fortran (runtime_error)... ✅ PASSED
  Expected: runtime_error
  Got:      runtime_error
  Response: {"status":"runtime_error","build":{"status":"ok","stdout":"","stderr":"","duration_ms":117},"tests":[{"status":"runtime_error","stdout":"","stderr":"ERROR STOP 3\n\nError termination. Backtrace:\n#0  0x777a3f71f8c2 in ???\n#1  0x777a3f7203b9 in ???\n#2  0x597e34a1616b in MAIN__\n#3  0x597e34a161a4 in main\n","duration_ms":19,"memory_peak_kb":256}]}
Testing fortran (time_exceeded)... ✅ PASSED
  Expected: time_exceeded
  Got:      time_exceeded
  Response: {"status":"time_exceeded","build":{"status":"ok","stdout":"","stderr":"","duration_ms":139},"tests":[{"status":"time_exceeded","stdout":"","stderr":"","duration_ms":1009,"memory_peak_kb":512}]}
Testing fortran (wrong_output)... ✅ PASSED
  Expected: wrong_output
  Got:      wrong_output
  Response: {"status":"wrong_output","build":{"status":"ok","stdout":"","stderr":"","duration_ms":104},"tests":[{"status":"wrong_output","stdout":"42\n","stderr":"","duration_ms":16,"memory_peak_kb":0}]}
Testing typescript (accepted)... ✅ PASSED
  Expected: accepted
  Got:      accepted
  Response: {"status":"accepted","build":{"status":"ok","stdout":"","stderr":"","duration_ms":3144},"tests":[{"status":"accepted","stdout":"42\n","stderr":"","duration_ms":186,"memory_peak_kb":15408}]}
Testing typescript (build_failed)... ✅ PASSED
  Expected: build_failed
  Got:      build_failed
  Response: {"status":"build_failed","build":{"status":"failed","stdout":"solution.ts(1,9): error TS1005: ',' expected.\n","stderr":"","duration_ms":3125},"tests":[{"status":"not_executed","stdout":"","stderr":"","duration_ms":0,"memory_peak_kb":0}]}
Testing typescript (runtime_error)... ✅ PASSED
  Expected: runtime_error
  Got:      runtime_error
  Response: {"status":"runtime_error","build":{"status":"ok","stdout":"","stderr":"","duration_ms":3089},"tests":[{"status":"runtime_error","stdout":"","stderr":"/sandbox/solution.js:2\nthrow new Error('boom');\n^\n\nError: boom\n    at Object.\u003canonymous\u003e (/sandbox/solution.js:2:7)\n    at Module._compile (node:internal/modules/cjs/loader:1364:14)\n    at Module._extensions..js (node:internal/modules/cjs/loader:1422:10)\n    at Module.load (node:internal/modules/cjs/loader:1203:32)\n    at Module._load (node:internal/modules/cjs/loader:1019:12)\n    at Function.executeUserEntryPoint [as runMain] (node:internal/modules/run_main:128:12)\n    at node:internal/main/run_main_module:28:49\n\nNode.js v18.20.4\n","duration_ms":174,"memory_peak_kb":14580}]}
Testing typescript (time_exceeded)... ✅ PASSED
  Expected: time_exceeded
  Got:      time_exceeded
  Response: {"status":"time_exceeded","build":{"status":"ok","stdout":"","stderr":"","duration_ms":3014},"tests":[{"status":"time_exceeded","stdout":"","stderr":"","duration_ms":1007,"memory_peak_kb":14032}]}
Testing typescript (wrong_output)... ✅ PASSED
  Expected: wrong_output
  Got:      wrong_output
  Response: {"status":"wrong_output","build":{"status":"ok","stdout":"","stderr":"","duration_ms":3056},"tests":[{"status":"wrong_output","stdout":"42\n","stderr":"","duration_ms":213,"memory_peak_kb":15324}]}

==============================
 Results: 14/14 passed
 ALL PAYLOADS PASSED
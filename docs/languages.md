# Language Registry

`goboxd` uses a YAML-based registry to define how different programming languages are handled. This allows adding new languages without changing Go code.

## Configuration File (`languages.yaml`)

Each language is defined as a block under the `languages` key.

### Fields
- `id`: Unique identifier (e.g., `py3`).
- `name`: Human-readable name.
- `source_filename`: File name used to save the source code.
- `artifact`: (Optional) File name for the compiled artifact.
- `build`: (Optional) Configuration for the compilation stage.
  - `cmd`: Binary to run (e.g., `/usr/bin/g++`).
  - `args`: Arguments with placeholders.
  - `flag_allowlist`: Glob patterns for allowed flags (e.g., `["-O*"]`).
  - `limits`: Wall time, memory, and process caps for the build.
- `run`: Configuration for the execution stage. Same structure as `build`.

### Placeholder System
Arguments in `cmd` use placeholders resolved at runtime:
- `{{source}}`: Path to the source file in the sandbox.
- `{{artifact}}`: Path to the compiled artifact in the sandbox.
- `{{flags}}`: Client-supplied flags (validated against allowlist).

## Currently Registered Languages

### Python 3 (`py3`)
- **Source**: `solution.py`
- **Build**: None
- **Run Limits**: 9s, 100MB, 100 processes
- **Run Command**: `/usr/bin/python3 {{source}}`

---

## Adding a New Language
To add a language (e.g., Rust):
1. Update `languages.yaml` with the new block.
2. Ensure the toolchain (e.g., `rustc`) is installed in the Docker image.
3. Restart the service.

No Go code changes are required for standard languages.

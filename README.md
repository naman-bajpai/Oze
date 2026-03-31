# oze

> AI-driven feature loop — describe a feature, let Claude implement it, run tests, repeat until green.

## How it works

1. You run `oze "Add rate limiting to /api/login"`
2. oze auto-detects your test command (or you pass `--test`)
3. oze sends a prompt to `claude --print --dangerously-skip-permissions`
4. oze runs the tests and captures the output
5. **Tests pass** → prints a success summary and exits 0
6. **Tests fail** → feeds the failure output back to Claude and retries
7. Repeats up to `--max` iterations (default: 10)

## Requirements

- Go 1.22+
- [Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code) installed and on `$PATH` as `claude`

## Installation

```bash
# Clone and build
git clone https://github.com/yourusername/oze.git
cd oze
go build -o oze .

# Or install globally
go install github.com/yourusername/oze@latest
```

> **Before pushing**: replace `yourusername` in `go.mod` with your real GitHub username:
> ```bash
> sed -i '' 's/yourusername/YOUR_REAL_USERNAME/g' go.mod main.go internal/cli/cli.go
> ```

## Usage

```bash
oze "Add rate limiting to /api/login"
oze --test "pytest -x" "Add input validation to signup form"
oze --dry-run "Refactor auth module"
oze --verbose --max 5 "Fix the pagination bug"
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--test <cmd>` | auto-detect | Override the test command |
| `--max <n>` | 10 | Max iterations before giving up |
| `--dry-run` | false | Print the Claude prompt, don't execute |
| `--verbose` | false | Stream Claude output live |
| `--no-color` | false | Disable ANSI colors |
| `--version` | — | Print version and exit |

## Test command auto-detection order

1. `CLAUDE.md` — line starting with `Test command:` or `test:`
2. `package.json` → `scripts.test` → `npm test`
3. `Makefile` with `test:` target → `make test`
4. `pytest.ini`, `setup.cfg`, or `pyproject.toml` → `pytest`
5. `Cargo.toml` → `cargo test`
6. `go.mod` → `go test ./...`
7. `Gemfile` + `Rakefile` → `bundle exec rake test`
8. `pom.xml` → `mvn test`
9. `build.gradle` → `./gradlew test`
10. None found → exit with error (use `--test`)

## License

MIT

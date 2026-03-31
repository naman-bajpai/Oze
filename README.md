


# oze

> A terminal CLI that implements a feature, runs your tests, fixes failures, and keeps going until it works — powered by Claude Code.

```
oze "Add rate limiting to the /api/login endpoint"
```

oze loops: **implement → test → fix → repeat** — up to 10 times, automatically. You describe what you want. Claude does the work. oze doesn't stop until tests pass.

---

## Install

### Prerequisites
- [Go 1.21+](https://go.dev/dl/)
- [Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code/overview) installed and authenticated (`npm install -g @anthropic-ai/claude-code`)

### Build from source

```bash
git clone https://github.com/yourusername/oze
cd oze
go build -o oze .
sudo mv oze /usr/local/bin/   # or anywhere on your PATH
```

### Verify

```bash
oze --help
```

---

## Usage

```bash
oze "<feature description>"
oze [flags] "<feature description>"
```

### Examples

```bash
# Basic — oze auto-detects the test command
oze "Add input validation to the signup form"

# Specify the test command explicitly
oze --test "npm test" "Implement JWT refresh token rotation"

# Use pytest with fail-fast, limit to 5 iterations
oze --test "pytest -x" --max 5 "Add rate limiting to /api/login"

# See what prompt oze would send without running anything
oze --dry-run "Refactor the auth module"

# Stream full Claude output
oze --verbose "Fix all TypeScript errors in src/api/"
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--test <cmd>` | auto-detect | Test command to run after each iteration |
| `--max <n>` | 10 | Max iterations before giving up |
| `--dry-run` | false | Print the Claude prompt, don't execute |
| `--verbose` | false | Stream full Claude output live |
| `--no-color` | false | Disable colored terminal output |

---

## How it works

```
oze "your feature"
       │
       ▼
  detect test command
  (or use --test flag)
       │
       ▼
  ┌─────────────────────────────────┐
  │  [iteration N]                  │
  │                                 │
  │  1. Build prompt for Claude     │
  │  2. Run: claude --print ...     │
  │  3. Run: <test command>         │
  │                                 │
  │  tests pass? ──► Done! ✔       │
  │  tests fail? ──► next iteration │
  │  N == max?  ──► report + exit 1 │
  └─────────────────────────────────┘
```

On each iteration oze tells Claude exactly what failed (the full test output) so it can make targeted fixes — not blind rewrites.

---

## Test command auto-detection

oze checks these in order until it finds a match:

| File | Command used |
|------|-------------|
| `CLAUDE.md` (line starting with `Test command:`) | whatever you wrote |
| `package.json` → `scripts.test` | `npm test` |
| `Makefile` with `test:` target | `make test` |
| `pytest.ini`, `setup.cfg`, `pyproject.toml` | `pytest` |
| `Cargo.toml` | `cargo test` |
| `go.mod` | `go test ./...` |
| `Gemfile` + `Rakefile` | `bundle exec rake test` |
| `Gemfile` | `bundle exec rspec` |
| `pom.xml` | `mvn test` |
| `build.gradle` / `build.gradle.kts` | `./gradlew test` |

If nothing matches, oze exits with a helpful error telling you to use `--test`.

**Tip:** add a line to your `CLAUDE.md` to always win detection:
```
Test command: npm run test:unit
```

---

## Sharing with your team

oze is a single binary — just share the repo. Anyone can clone and `go build` it:

```bash
git clone https://github.com/yourusername/oze
cd oze && go build -o oze . && sudo mv oze /usr/local/bin/
```

Or add a `Makefile` target to your project:
```makefile
install-oze:
    go install github.com/yourusername/oze@latest
```

---

## Tips

- Run `oze --dry-run "..."` first on complex tasks to preview what Claude will receive.
- Add `Test command:` to your `CLAUDE.md` to skip detection entirely.
- Use `--max 20` for larger features that need more iterations.
- oze exits with code `0` on success, `1` on failure — pipe it into CI if you like.
# Oze

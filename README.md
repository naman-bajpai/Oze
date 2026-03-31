# oze

> AI-driven feature loop — describe a feature, let Claude implement it, run tests, repeat until green.

## How it works

1. You run `oze "Add rate limiting to /api/login"`
2. oze auto-detects your test command (or you pass `--test` / `--no-test`)
3. oze sends a focused prompt to `claude --print --dangerously-skip-permissions`
4. oze runs the tests and captures the output
5. **Tests pass** → prints a success summary and exits 0
6. **Tests fail** → feeds the failure output back to Claude and retries
7. Repeats up to `--max` iterations (default: 10)

## Requirements

- Go 1.22+
- [Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code) installed and on `$PATH` as `claude`

## Installation

```bash
git clone https://github.com/naman-bajpai/oze.git
cd oze
make install   # builds and copies to $GOPATH/bin
```

### Update after pulling new changes

```bash
git pull
make install
```

## Usage

```bash
# Basic
oze "Add rate limiting to /api/login"

# With specialist context
oze --frontend "Add dark mode toggle to the nav bar"
oze --backend --test "go test ./..." "Add JWT refresh endpoint"
oze --security "Sanitise user input on the signup form"

# Projects without automated tests
oze --no-test "Rename component to EchoSpace"
oze --no-test --frontend --model haiku "Fix spacing on the dashboard"

# Other
oze --dry-run "Refactor auth module"
oze --verbose --max 5 "Fix the pagination bug"
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--test <cmd>` | auto-detect | Override the test command |
| `--no-test` | false | Skip tests — run Claude once and exit |
| `--model <name>` | default | Claude model: `haiku`, `sonnet`, `opus` |
| `--max <n>` | 10 | Max iterations before giving up |
| `--dry-run` | false | Print the Claude prompt, don't execute |
| `--verbose` | false | Stream Claude output live |
| `--no-color` | false | Disable ANSI colors |
| `--version` | — | Print version and exit |

## Specialist modes

Activates a focused system prompt so Claude stays in the right layer of the stack.
Only one specialist flag may be used at a time.

| Flag | Focus |
|------|-------|
| `--frontend` | React, TypeScript, Tailwind CSS, accessibility, Core Web Vitals |
| `--backend` | APIs, auth, input validation, error handling, security |
| `--mobile` | React Native, Expo, iOS/Android platform UX |
| `--database` | Schema design, indexes, query optimisation, migrations |
| `--devops` | CI/CD, Docker, infra-as-code, secrets management |
| `--security` | OWASP Top 10, auth/authz, secure defaults |

## Test command auto-detection order

1. `CLAUDE.md` — line starting with `Test command:` or `test:`
2. `package.json` → `scripts.test` → `pnpm test` / `yarn test` / `npm test`
3. `Makefile` with `test:` target → `make test`
4. `pytest.ini`, `setup.cfg`, or `pyproject.toml` → `pytest`
5. `Cargo.toml` → `cargo test`
6. `go.mod` → `go test ./...`
7. `Gemfile` + `Rakefile` → `bundle exec rake test`
8. `pom.xml` → `mvn test`
9. `build.gradle` → `./gradlew test`
10. None found → error (use `--test` or `--no-test`)

## License

MIT

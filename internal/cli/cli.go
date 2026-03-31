// Package cli implements flag parsing and the main oze loop.
package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/naman-bajpai/oze/internal/claude"
	"github.com/naman-bajpai/oze/internal/detector"
	"github.com/naman-bajpai/oze/internal/logger"
	"github.com/naman-bajpai/oze/internal/runner"
	"github.com/naman-bajpai/oze/internal/specialist"
)

const version = "0.1.0"

// Run is the entry point called from main.go.
func Run() {
	// ── Flag definitions ──────────────────────────────────────────────────
	testCmd    := flag.String("test", "", "Override auto-detected test command")
	noTest     := flag.Bool("no-test", false, "Skip test loop — run Claude once and exit")
	model      := flag.String("model", "", "Claude model to use (e.g. haiku, sonnet, opus)")
	maxIter    := flag.Int("max", 10, "Max iterations before giving up (default 10)")
	dryRun     := flag.Bool("dry-run", false, "Print the Claude prompt without executing")
	verbose    := flag.Bool("verbose", false, "Stream Claude output live to the terminal")
	noColor    := flag.Bool("no-color", false, "Disable ANSI colors")
	showVersion := flag.Bool("version", false, "Print version and exit")

	// Specialist flags
	isFrontend := flag.Bool("frontend", false, "Activate frontend specialist (React, TypeScript, Tailwind, a11y)")
	isBackend  := flag.Bool("backend", false, "Activate backend specialist (APIs, auth, security, DB)")
	isMobile   := flag.Bool("mobile", false, "Activate mobile specialist (React Native, Expo, iOS/Android)")
	isDatabase := flag.Bool("database", false, "Activate database specialist (schema, indexes, migrations)")
	isDevOps   := flag.Bool("devops", false, "Activate DevOps specialist (CI/CD, Docker, infra-as-code)")
	isSecurity := flag.Bool("security", false, "Activate security specialist (OWASP, auth, secrets)")

	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Printf("oze version %s\n", version)
		os.Exit(0)
	}

	// ── Positional arg: feature description ───────────────────────────────
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: feature description is required.")
		fmt.Fprintln(os.Stderr)
		usage()
		os.Exit(1)
	}
	feature := flag.Arg(0)

	log := logger.New(*noColor)

	// ── Resolve specialist ─────────────────────────────────────────────────
	var specialistPrompt string
	specialistMap := map[*bool]specialist.Role{
		isFrontend: specialist.Frontend,
		isBackend:  specialist.Backend,
		isMobile:   specialist.Mobile,
		isDatabase: specialist.Database,
		isDevOps:   specialist.DevOps,
		isSecurity: specialist.Security,
	}
	activeCount := 0
	for flagPtr, role := range specialistMap {
		if *flagPtr {
			activeCount++
			if activeCount > 1 {
				fmt.Fprintln(os.Stderr, "Error: only one specialist flag may be used at a time.")
				os.Exit(1)
			}
			p, _ := specialist.Prompt(role)
			specialistPrompt = p
			log.Info(fmt.Sprintf("Specialist: %s", role))
		}
	}

	// ── Resolve working directory ──────────────────────────────────────────
	workDir, err := filepath.Abs(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving working directory: %v\n", err)
		os.Exit(1)
	}

	// ── Auto-detect test command ───────────────────────────────────────────
	cmd := *testCmd
	if *noTest {
		cmd = ""
	} else if cmd == "" {
		detected, err := detector.Detect(workDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			fmt.Fprintln(os.Stderr, "Tip: use --no-test to run Claude once without a test loop.")
			os.Exit(1)
		}
		cmd = detected
		log.Info(fmt.Sprintf("Auto-detected test command: %s", cmd))
	} else {
		log.Info(fmt.Sprintf("Using test command: %s", cmd))
	}

	// ── Dry-run mode ───────────────────────────────────────────────────────
	if *dryRun {
		prompt := claude.BuildPrompt(1, feature, cmd, "")
		log.DryRun(prompt)
		os.Exit(0)
	}

	// ── Main loop ─────────────────────────────────────────────────────────
	log.Banner(feature)

	// --no-test: run Claude once and exit without a test loop
	if *noTest {
		log.Iteration(1, 1, fmt.Sprintf("Calling Claude to implement: %s", feature))
		prompt := claude.BuildPrompt(1, feature, "", "")
		_, err := claude.Run(prompt, claude.Options{
			Verbose:      *verbose,
			WorkDir:      workDir,
			Model:        *model,
			SystemPrompt: specialistPrompt,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "[oze] Claude error: %v\n", err)
			os.Exit(1)
		}
		log.Success(feature, 1)
		os.Exit(0)
	}

	var lastTestOutput string

	for i := 1; i <= *maxIter; i++ {
		log.Iteration(i, *maxIter, fmt.Sprintf("Calling Claude to implement: %s", feature))

		prompt := claude.BuildPrompt(i, feature, cmd, lastTestOutput)

		_, err := claude.Run(prompt, claude.Options{
			Verbose:      *verbose,
			WorkDir:      workDir,
			Model:        *model,
			SystemPrompt: specialistPrompt,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "[oze] Claude error on iteration %d: %v\n", i, err)
			// Non-fatal: continue to run tests; Claude may have made partial changes.
		}

		log.Info(fmt.Sprintf("Running tests: %s", cmd))

		result, err := runner.Run(workDir, cmd, 300)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[oze] Failed to execute test command: %v\n", err)
			os.Exit(1)
		}

		if result.Passed {
			log.TestPass()
			log.Success(feature, i)
			os.Exit(0)
		}

		lastTestOutput = result.Output
		log.TestFail(result.Output)

		if i == *maxIter {
			log.MaxReached(*maxIter, lastTestOutput)
			os.Exit(1)
		}

		log.Info(fmt.Sprintf("Tests failed — feeding output back to Claude (iteration %d/%d)", i+1, *maxIter))
	}
}

// usage prints the help text.
func usage() {
	fmt.Fprintf(os.Stderr, `oze — AI-driven feature loop (v%s)

USAGE
  oze [flags] "feature description"

EXAMPLES
  oze "Add rate limiting to /api/login"
  oze --frontend "Add dark mode toggle to the nav bar"
  oze --backend --test "go test ./..." "Add JWT refresh endpoint"
  oze --no-test --model haiku "Rename component to EchoSpace"
  oze --dry-run --security "Sanitise user input on signup form"

SPECIALISTS
  --frontend   React, TypeScript, Tailwind, accessibility, Core Web Vitals
  --backend    APIs, auth, input validation, error handling, security
  --mobile     React Native, Expo, iOS/Android platform UX
  --database   Schema design, indexes, query optimisation, migrations
  --devops     CI/CD, Docker, infra-as-code, secrets management
  --security   OWASP Top 10, auth/authz, secure defaults

FLAGS
`, version)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
NOTES
  oze auto-detects the test command from your project files in this order:
    CLAUDE.md > package.json > Makefile > pytest.ini/setup.cfg/pyproject.toml >
    Cargo.toml > go.mod > Gemfile+Rakefile > pom.xml > build.gradle

  Use --test to override if auto-detection picks the wrong command.

  Module path placeholder: replace "yourusername" in go.mod with your GitHub handle.
`)
}

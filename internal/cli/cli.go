package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/oze/internal/claude"
	"github.com/yourusername/oze/internal/detector"
	"github.com/yourusername/oze/internal/logger"
	"github.com/yourusername/oze/internal/runner"
)

const maxIterations = 10

const usage = `oze — iterative feature completion loop powered by Claude Code

Usage:
  oze [flags] "<feature description>"

Examples:
  oze "Add input validation to the signup form"
  oze --test "npm test" "Implement JWT refresh token rotation"
  oze --test "pytest -x" --max 5 "Add rate limiting to /api/login"
  oze --dry-run "Refactor the auth module"

Flags:
`

// Config holds all resolved runtime options.
type Config struct {
	Feature     string
	TestCmd     string
	MaxIter     int
	DryRun      bool
	Verbose     bool
	NoColor     bool
}

func Run(args []string) error {
	fs := flag.NewFlagSet("oze", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		fs.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}

	testCmd  := fs.String("test", "", "Test command to run (overrides auto-detection)")
	maxIter  := fs.Int("max", maxIterations, "Maximum number of implement→test iterations")
	dryRun   := fs.Bool("dry-run", false, "Print the Claude prompt without executing")
	verbose  := fs.Bool("verbose", false, "Stream full Claude output")
	noColor  := fs.Bool("no-color", false, "Disable colored output")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return fmt.Errorf("feature description required")
	}

	feature := strings.Join(fs.Args(), " ")

	cfg := Config{
		Feature: feature,
		TestCmd: *testCmd,
		MaxIter: *maxIter,
		DryRun:  *dryRun,
		Verbose: *verbose,
		NoColor: *noColor,
	}

	return runLoop(cfg)
}

func runLoop(cfg Config) error {
	log := logger.New(!cfg.NoColor)

	// Resolve test command
	testCmd := cfg.TestCmd
	if testCmd == "" {
		detected, err := detector.FindTestCommand(".")
		if err != nil || detected == "" {
			log.Warn("Could not auto-detect test command. Use --test to specify one.")
			log.Warn("Example: oze --test \"npm test\" \"" + cfg.Feature + "\"")
			return fmt.Errorf("no test command found")
		}
		testCmd = detected
		log.Info(fmt.Sprintf("Detected test command: %s", testCmd))
	} else {
		log.Info(fmt.Sprintf("Using test command: %s", testCmd))
	}

	log.Banner(fmt.Sprintf("oze — %s", cfg.Feature))

	if cfg.DryRun {
		prompt := claude.BuildPrompt(cfg.Feature, testCmd, 1, cfg.MaxIter, "")
		fmt.Println("\n--- Claude prompt (dry run) ---")
		fmt.Println(prompt)
		fmt.Println("--- end ---")
		return nil
	}

	// Check claude CLI is available
	if err := claude.CheckAvailable(); err != nil {
		return err
	}

	var lastOutput string

	for i := 1; i <= cfg.MaxIter; i++ {
		log.Step(i, cfg.MaxIter, "Calling Claude")

		prompt := claude.BuildPrompt(cfg.Feature, testCmd, i, cfg.MaxIter, lastOutput)

		claudeOut, err := claude.Run(prompt, cfg.Verbose)
		if err != nil {
			log.Error(fmt.Sprintf("Claude failed on iteration %d: %v", i, err))
			return err
		}

		log.Step(i, cfg.MaxIter, fmt.Sprintf("Running tests: %s", testCmd))

		testOut, testErr := runner.Run(testCmd)
		lastOutput = testOut

		if testErr == nil {
			log.Success(fmt.Sprintf("All tests passed on iteration %d!", i))
			log.Box("Test output", trimOutput(testOut, 40))
			log.Done(cfg.Feature, i)
			return nil
		}

		log.Fail(fmt.Sprintf("Tests failed (iteration %d/%d)", i, cfg.MaxIter))
		if cfg.Verbose {
			log.Box("Test output", trimOutput(testOut, 30))
		} else {
			log.Box("Test output (last 15 lines)", trimOutput(testOut, 15))
		}

		if i == cfg.MaxIter {
			log.Error(fmt.Sprintf("Reached max iterations (%d). Feature not complete.", cfg.MaxIter))
			log.Warn("Try: oze --max 20 \"" + cfg.Feature + "\"")
			return fmt.Errorf("max iterations reached without passing tests")
		}

		log.Info(fmt.Sprintf("Retrying... (%d/%d)", i, cfg.MaxIter))
		_ = claudeOut // claudeOut already streamed/printed if verbose
	}

	return nil
}

func trimOutput(s string, maxLines int) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) <= maxLines {
		return strings.TrimSpace(s)
	}
	kept := lines[len(lines)-maxLines:]
	return fmt.Sprintf("... (%d lines omitted) ...\n%s", len(lines)-maxLines, strings.Join(kept, "\n"))
}

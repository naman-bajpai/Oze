// Package logger provides colored terminal output for oze.
package logger

import (
	"fmt"
	"strings"
	"time"
)

// ANSI color codes
const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	gray   = "\033[90m"
)

// Logger handles all terminal output with optional color support.
type Logger struct {
	noColor bool
}

// New creates a new Logger. Pass noColor=true to disable ANSI codes.
func New(noColor bool) *Logger {
	return &Logger{noColor: noColor}
}

func (l *Logger) color(code string) string {
	if l.noColor {
		return ""
	}
	return code
}

// Banner prints the startup banner with the feature description.
func (l *Logger) Banner(feature string) {
	width := 60
	if len(feature)+6 > width {
		width = len(feature) + 6
	}
	border := strings.Repeat("─", width)
	fmt.Printf("\n%s%s╭%s╮%s\n", l.color(cyan), l.color(bold), border, l.color(reset))
	title := "  oze — AI-driven feature loop  "
	pad := strings.Repeat(" ", (width-len(title))/2)
	fmt.Printf("%s%s│%s%s%s%s│%s\n",
		l.color(cyan), l.color(bold),
		l.color(reset), pad+title+pad,
		l.color(cyan), l.color(bold), l.color(reset))

	// Feature line
	featureLabel := "  Feature: " + feature
	featurePad := strings.Repeat(" ", width-len(featureLabel))
	fmt.Printf("%s%s│%s %s%s%s%s│%s\n",
		l.color(cyan), l.color(bold),
		l.color(reset),
		l.color(yellow), featureLabel[1:], featurePad,
		l.color(cyan)+l.color(bold), l.color(reset))

	fmt.Printf("%s%s╰%s╯%s\n\n", l.color(cyan), l.color(bold), border, l.color(reset))
}

// Iteration prints the start of an iteration.
func (l *Logger) Iteration(n, max int, intent string) {
	ts := time.Now().Format("15:04:05")
	fmt.Printf("%s%s%s %s[%d/%d]%s %s\n",
		l.color(gray), ts, l.color(reset),
		l.color(cyan)+l.color(bold), n, max, l.color(reset),
		intent)
}

// Info prints a general info line.
func (l *Logger) Info(msg string) {
	fmt.Printf("  %s→%s %s\n", l.color(cyan), l.color(reset), msg)
}

// TestPass prints a test-pass message.
func (l *Logger) TestPass() {
	fmt.Printf("  %s%s✔ Tests passed%s\n", l.color(green), l.color(bold), l.color(reset))
}

// TestFail prints a test-fail message and the last lines of output in a box.
func (l *Logger) TestFail(output string) {
	fmt.Printf("  %s%s✖ Tests failed%s\n", l.color(red), l.color(bold), l.color(reset))

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	tail := lines
	if len(lines) > 15 {
		tail = lines[len(lines)-15:]
	}

	width := 60
	for _, ln := range tail {
		if len(ln) > width {
			width = len(ln) + 4
		}
	}
	border := strings.Repeat("─", width)
	fmt.Printf("  %s┌%s┐%s\n", l.color(red), border, l.color(reset))
	for _, ln := range tail {
		pad := strings.Repeat(" ", width-len(ln))
		fmt.Printf("  %s│%s %s%s %s│%s\n",
			l.color(red), l.color(reset),
			ln, pad,
			l.color(red), l.color(reset))
	}
	fmt.Printf("  %s└%s┘%s\n\n", l.color(red), border, l.color(reset))
}

// DryRun prints the prompt that would be sent to Claude.
func (l *Logger) DryRun(prompt string) {
	fmt.Printf("%s%s── DRY RUN: Claude prompt ──%s\n", l.color(bold), l.color(yellow), l.color(reset))
	fmt.Println(prompt)
	fmt.Printf("%s%s── END DRY RUN ──%s\n\n", l.color(bold), l.color(yellow), l.color(reset))
}

// Success prints the final success summary.
func (l *Logger) Success(feature string, iterations int) {
	width := 60
	if len(feature)+14 > width {
		width = len(feature) + 14
	}
	border := strings.Repeat("─", width)
	fmt.Printf("\n%s%s╭%s╮%s\n", l.color(green), l.color(bold), border, l.color(reset))

	title := "  ✔ Feature implemented successfully  "
	titlePad := strings.Repeat(" ", width-len(title))
	fmt.Printf("%s%s│%s%s%s%s%s│%s\n",
		l.color(green), l.color(bold),
		l.color(reset), l.color(bold), title, titlePad,
		l.color(green)+l.color(bold), l.color(reset))

	featureLine := fmt.Sprintf("  Feature:    %s", feature)
	featurePad := strings.Repeat(" ", width-len(featureLine))
	fmt.Printf("%s%s│%s%s%s%s│%s\n",
		l.color(green), l.color(bold),
		l.color(reset), featureLine, featurePad,
		l.color(green)+l.color(bold), l.color(reset))

	iterLine := fmt.Sprintf("  Iterations: %d", iterations)
	iterPad := strings.Repeat(" ", width-len(iterLine))
	fmt.Printf("%s%s│%s%s%s%s│%s\n",
		l.color(green), l.color(bold),
		l.color(reset), iterLine, iterPad,
		l.color(green)+l.color(bold), l.color(reset))

	fmt.Printf("%s%s╰%s╯%s\n\n", l.color(green), l.color(bold), border, l.color(reset))
}

// MaxReached prints the status when max iterations are exhausted.
func (l *Logger) MaxReached(max int, lastOutput string) {
	fmt.Printf("\n%s%s✖ Max iterations (%d) reached without passing tests.%s\n",
		l.color(red), l.color(bold), max, l.color(reset))
	fmt.Printf("%sLast test output:%s\n", l.color(yellow), l.color(reset))
	l.TestFail(lastOutput)
}

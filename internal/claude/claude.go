// Package claude invokes the Claude Code CLI and captures its output.
package claude

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Options controls how the Claude CLI is invoked.
type Options struct {
	Verbose    bool   // stream Claude output live to the terminal
	WorkDir    string // working directory for the subprocess
	Model      string // model override (e.g. "haiku", "sonnet")
	SystemPrompt string // extra system prompt appended via --append-system-prompt
}

// Run sends prompt to `claude --print --dangerously-skip-permissions` and
// returns the combined output. If opts.Verbose is true the output is also
// streamed live to os.Stdout.
func Run(prompt string, opts Options) (string, error) {
	args := []string{
		"--print",
		"--dangerously-skip-permissions",
		"--no-session-persistence",
	}
	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}
	if opts.SystemPrompt != "" {
		args = append(args, "--append-system-prompt", opts.SystemPrompt)
	}

	cmd := exec.CommandContext(context.Background(), "claude", args...)
	cmd.Dir = opts.WorkDir

	// Feed the prompt via stdin
	cmd.Stdin = strings.NewReader(prompt)

	var buf bytes.Buffer

	if opts.Verbose {
		// Tee: write to both buffer and stdout
		cmd.Stdout = io.MultiWriter(&buf, os.Stdout)
		cmd.Stderr = io.MultiWriter(&buf, os.Stderr)
	} else {
		cmd.Stdout = &buf
		cmd.Stderr = &buf
	}

	if err := cmd.Run(); err != nil {
		combined := strings.TrimRight(buf.String(), "\n")
		return combined, fmt.Errorf("claude exited with error: %w\noutput:\n%s", err, combined)
	}

	return strings.TrimRight(buf.String(), "\n"), nil
}

// maxFailureLines is the maximum number of lines from test output fed back to Claude.
const maxFailureLines = 60

// BuildPrompt constructs the prompt string to send to Claude.
//
// iteration is 1-based.
// feature is the user-supplied feature description.
// testCmd is the detected/provided test command (empty when --no-test).
// prevFailure is the test output from the previous iteration (empty on iter 1).
func BuildPrompt(iteration int, feature, testCmd, prevFailure string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Implement: %s\n", feature))
	if testCmd != "" {
		sb.WriteString(fmt.Sprintf("Test cmd (do NOT run): %s\n", testCmd))
	}
	sb.WriteString("Make minimal targeted changes only.\n")

	if prevFailure != "" {
		sb.WriteString("\nPrevious test failure — fix this:\n```\n")
		sb.WriteString(truncateLines(prevFailure, maxFailureLines))
		sb.WriteString("\n```\n")
	}

	if iteration > 1 {
		sb.WriteString(fmt.Sprintf("\n(Attempt %d)\n", iteration))
	}

	return sb.String()
}

// truncateLines returns the last n lines of s (or all of s if fewer than n lines).
func truncateLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) <= n {
		return strings.Join(lines, "\n")
	}
	omitted := len(lines) - n
	tail := lines[omitted:]
	return fmt.Sprintf("[...%d lines omitted...]\n", omitted) + strings.Join(tail, "\n")
}

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
	Verbose bool   // stream Claude output live to the terminal
	WorkDir string // working directory for the subprocess
}

// Run sends prompt to `claude --print --dangerously-skip-permissions` and
// returns the combined output. If opts.Verbose is true the output is also
// streamed live to os.Stdout.
func Run(prompt string, opts Options) (string, error) {
	args := []string{
		"--print",
		"--dangerously-skip-permissions",
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

// BuildPrompt constructs the prompt string to send to Claude.
//
// iteration is 1-based.
// feature is the user-supplied feature description.
// testCmd is the detected/provided test command.
// prevFailure is the test output from the previous iteration (empty on iter 1).
func BuildPrompt(iteration int, feature, testCmd, prevFailure string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[oze] Iteration %d — implement: %s\n\n", iteration, feature))

	sb.WriteString("## Task\n")
	sb.WriteString(fmt.Sprintf("You are helping implement the following feature:\n\n  %s\n\n", feature))

	sb.WriteString("## Instructions\n")
	sb.WriteString("1. Read the relevant source files in the current directory to understand the codebase.\n")
	sb.WriteString("2. Implement the feature described above by editing or creating files as needed.\n")
	sb.WriteString(fmt.Sprintf("3. The test command for this project is: `%s`\n", testCmd))
	sb.WriteString("4. Do NOT run the tests yourself — the caller will run them after you finish.\n")
	sb.WriteString("5. Make minimal, targeted changes. Do not refactor unrelated code.\n")
	sb.WriteString("6. When you have finished all changes, output exactly the word DONE on its own line.\n")

	if prevFailure != "" {
		sb.WriteString("\n## Previous test failure\n")
		sb.WriteString("The tests ran after your last set of changes and FAILED. Here is the output:\n\n")
		sb.WriteString("```\n")
		sb.WriteString(prevFailure)
		sb.WriteString("\n```\n\n")
		sb.WriteString("Analyse the failure carefully and make a targeted fix before outputting DONE.\n")
	}

	return sb.String()
}

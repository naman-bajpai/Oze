package claude

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CheckAvailable verifies that the claude CLI is on PATH.
func CheckAvailable() error {
	_, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf(`claude CLI not found on PATH.

Install it with:
  npm install -g @anthropic-ai/claude-code

Then run: claude --version to verify.`)
	}
	return nil
}

// BuildPrompt constructs the full prompt sent to Claude on each iteration.
func BuildPrompt(feature, testCmd string, iteration, maxIter int, lastTestOutput string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`You are running iteration %d of %d in an automated feature-completion loop called oze.

## Your task
%s

## Test command
%s

## What you must do this iteration
`, iteration, maxIter, feature, testCmd))

	if iteration == 1 {
		sb.WriteString(`This is the first iteration. 

1. Read the relevant source files and understand the codebase structure.
2. Implement the feature described above.
3. Make sure your changes are complete enough that the test command could plausibly pass.
4. Keep changes focused — only touch files relevant to this feature.
`)
	} else {
		sb.WriteString(fmt.Sprintf(`The previous iteration's test run FAILED. Here is the test output:

--- PREVIOUS TEST OUTPUT ---
%s
--- END TEST OUTPUT ---

1. Analyse the failure above carefully.
2. Identify the root cause — do not guess, read the error.
3. Make a targeted fix. Prefer small changes over rewrites.
4. Do not re-implement things that are already working.
`, lastTestOutput))
	}

	sb.WriteString(fmt.Sprintf(`
## Rules
- Do NOT run the test command yourself — oze handles that after you finish.
- Do NOT ask clarifying questions — make your best judgement and proceed.
- Keep changes minimal and focused.
- Log your intent in one line at the start of your response: "[oze] Iteration %d — <what you are doing>"
- When you are done making changes, say "DONE" on its own line.
`, iteration))

	return sb.String()
}

// Run executes: claude --print "<prompt>" and returns combined stdout+stderr.
// If verbose is true it also streams output to the terminal in real time.
func Run(prompt string, verbose bool) (string, error) {
	// Write prompt to a temp file to avoid shell quoting issues with long prompts
	tmp, err := os.CreateTemp("", "oze-prompt-*.txt")
	if err != nil {
		return "", fmt.Errorf("creating temp prompt file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(prompt); err != nil {
		return "", fmt.Errorf("writing prompt: %w", err)
	}
	tmp.Close()

	// claude --print reads from stdin or accepts a prompt arg.
	// We use: claude --print "$(cat file)" pattern via exec with stdin pipe.
	cmd := exec.Command("claude", "--print", "--dangerously-skip-permissions")
	cmd.Stdin = strings.NewReader(prompt)

	var buf bytes.Buffer
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = &buf
		cmd.Stderr = &buf
	}

	if err := cmd.Run(); err != nil {
		output := buf.String()
		return output, fmt.Errorf("claude exited with error: %w\n%s", err, output)
	}

	return buf.String(), nil
}

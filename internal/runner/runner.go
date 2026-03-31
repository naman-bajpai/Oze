// Package runner executes a shell test command and captures its output.
package runner

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

// Result holds the outcome of a test run.
type Result struct {
	Passed   bool
	Output   string // combined stdout + stderr
	Duration time.Duration
}

// Run executes cmdStr (a shell command) in dir and returns a Result.
// A non-zero exit code is treated as a test failure, not as an error.
func Run(dir, cmdStr string, timeoutSeconds int) (*Result, error) {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 300 // 5-minute default
	}
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeoutSeconds)*time.Second,
	)
	defer cancel()

	// Run through the shell so the user can pass things like "npm test -- --ci"
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
	cmd.Dir = dir

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	start := time.Now()
	runErr := cmd.Run()
	duration := time.Since(start)

	output := buf.String()
	output = strings.TrimRight(output, "\n")

	passed := runErr == nil
	// context.DeadlineExceeded means timeout — treat as failure, not error
	if ctx.Err() == context.DeadlineExceeded {
		output += "\n[oze] Test command timed out after " +
			time.Duration(timeoutSeconds).String() + "s"
		passed = false
		runErr = nil
	}

	return &Result{
		Passed:   passed,
		Output:   output,
		Duration: duration,
	}, runErr
}

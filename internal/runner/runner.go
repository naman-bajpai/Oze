package runner

import (
	"bytes"
	"os/exec"
	"strings"
)

// Run executes a shell command string and returns combined stdout+stderr.
// Returns a non-nil error if the command exits non-zero.
func Run(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return "", nil
	}

	c := exec.Command(parts[0], parts[1:]...)

	var buf bytes.Buffer
	c.Stdout = &buf
	c.Stderr = &buf

	err := c.Run()
	return buf.String(), err
}

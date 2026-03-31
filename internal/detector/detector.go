// Package detector auto-detects the test command for a project.
package detector

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Detect walks the given directory and returns the appropriate test command.
// It checks project files in priority order as specified by the oze spec.
func Detect(dir string) (string, error) {
	// 1. CLAUDE.md — line starting with "Test command:" or "test:"
	if cmd, ok := fromClaudeMD(dir); ok {
		return cmd, nil
	}

	// 2. package.json → scripts.test → prefer pnpm/yarn/npm based on lockfile
	if fileExists(dir, "package.json") {
		if hasNPMTestScript(dir) {
			if fileExists(dir, "pnpm-lock.yaml") || fileExists(dir, "pnpm-workspace.yaml") {
				return "pnpm test", nil
			}
			if fileExists(dir, "yarn.lock") {
				return "yarn test", nil
			}
			return "npm test", nil
		}
	}

	// 3. Makefile with "test:" target → "make test"
	if fileExists(dir, "Makefile") {
		if makefileHasTestTarget(dir) {
			return "make test", nil
		}
	}

	// 4. pytest.ini, setup.cfg, or pyproject.toml → "pytest"
	if fileExists(dir, "pytest.ini") || fileExists(dir, "setup.cfg") || fileExists(dir, "pyproject.toml") {
		return "pytest", nil
	}

	// 5. Cargo.toml → "cargo test"
	if fileExists(dir, "Cargo.toml") {
		return "cargo test", nil
	}

	// 6. go.mod → "go test ./..."
	if fileExists(dir, "go.mod") {
		return "go test ./...", nil
	}

	// 7. Gemfile + Rakefile → "bundle exec rake test"
	if fileExists(dir, "Gemfile") && fileExists(dir, "Rakefile") {
		return "bundle exec rake test", nil
	}

	// 8. pom.xml → "mvn test"
	if fileExists(dir, "pom.xml") {
		return "mvn test", nil
	}

	// 9. build.gradle → "./gradlew test"
	if fileExists(dir, "build.gradle") {
		return "./gradlew test", nil
	}

	// 10. Nothing found
	return "", fmt.Errorf(
		"could not auto-detect a test command in %q.\n"+
			"Please specify one with the --test flag (e.g. --test \"go test ./...\")",
		dir,
	)
}

// fileExists returns true if name exists inside dir.
func fileExists(dir, name string) bool {
	_, err := os.Stat(filepath.Join(dir, name))
	return err == nil
}

// fromClaudeMD parses CLAUDE.md looking for a "Test command:" or "test:" line.
func fromClaudeMD(dir string) (string, bool) {
	path := filepath.Join(dir, "CLAUDE.md")
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lower := strings.ToLower(line)
		for _, prefix := range []string{"test command:", "test:"} {
			if strings.HasPrefix(lower, prefix) {
				cmd := strings.TrimSpace(line[len(prefix):])
				if cmd != "" {
					return cmd, true
				}
			}
		}
	}
	return "", false
}

// hasNPMTestScript returns true if package.json defines a "test" script.
func hasNPMTestScript(dir string) bool {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return false
	}
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		// Even if we can't parse it, assume npm test exists if the file is present.
		return true
	}
	_, ok := pkg.Scripts["test"]
	return ok
}

// makefileHasTestTarget returns true if the Makefile contains a "test:" target.
func makefileHasTestTarget(dir string) bool {
	f, err := os.Open(filepath.Join(dir, "Makefile"))
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// A target line starts at column 0 and ends with ':'
		if strings.HasPrefix(line, "test:") || line == "test:" {
			return true
		}
	}
	return false
}

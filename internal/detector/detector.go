package detector

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

// FindTestCommand searches the given directory (usually ".") for a recognisable
// test command, in priority order. Returns ("", nil) if nothing is found.
func FindTestCommand(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	// 1. CLAUDE.md — look for lines like "Test command: npm test"
	if cmd := scanClaudeMD(abs); cmd != "" {
		return cmd, nil
	}

	// 2. package.json → scripts.test
	if cmd := scanPackageJSON(abs); cmd != "" {
		return cmd, nil
	}

	// 3. Makefile with a test target
	if hasFileWithTarget(abs, "Makefile", "test:") {
		return "make test", nil
	}

	// 4. Python — pytest markers
	for _, f := range []string{"pytest.ini", "setup.cfg", "pyproject.toml"} {
		if fileExists(filepath.Join(abs, f)) {
			return "pytest", nil
		}
	}

	// 5. Rust
	if fileExists(filepath.Join(abs, "Cargo.toml")) {
		return "cargo test", nil
	}

	// 6. Go
	if fileExists(filepath.Join(abs, "go.mod")) {
		return "go test ./...", nil
	}

	// 7. Ruby
	if fileExists(filepath.Join(abs, "Gemfile")) {
		if fileExists(filepath.Join(abs, "Rakefile")) {
			return "bundle exec rake test", nil
		}
		return "bundle exec rspec", nil
	}

	// 8. Java / Maven
	if fileExists(filepath.Join(abs, "pom.xml")) {
		return "mvn test", nil
	}

	// 9. Java / Gradle
	if fileExists(filepath.Join(abs, "build.gradle")) || fileExists(filepath.Join(abs, "build.gradle.kts")) {
		return "./gradlew test", nil
	}

	return "", nil
}

func scanClaudeMD(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		lower := strings.ToLower(strings.TrimSpace(line))
		for _, prefix := range []string{"test command:", "test:", "run tests:"} {
			if strings.HasPrefix(lower, prefix) {
				cmd := strings.TrimSpace(line[len(prefix):])
				if cmd != "" {
					return cmd
				}
			}
		}
	}
	return ""
}

func scanPackageJSON(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return ""
	}
	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return ""
	}
	if cmd, ok := pkg.Scripts["test"]; ok && cmd != "" {
		return "npm test"
	}
	return ""
}

func hasFileWithTarget(dir, filename, target string) bool {
	data, err := os.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		return false
	}
	return strings.Contains(string(data), target)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

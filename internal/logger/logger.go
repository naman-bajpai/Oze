package logger

import (
	"fmt"
	"strings"
	"time"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	cyan   = "\033[36m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	gray   = "\033[90m"
	blue   = "\033[34m"
)

// Logger handles colored terminal output.
type Logger struct {
	color bool
}

// New creates a Logger. Pass color=false to disable ANSI codes.
func New(color bool) *Logger {
	return &Logger{color: color}
}

func (l *Logger) c(code, s string) string {
	if !l.color {
		return s
	}
	return code + s + reset
}

// Banner prints the oze header.
func (l *Logger) Banner(title string) {
	line := strings.Repeat("─", 60)
	fmt.Println()
	fmt.Println(l.c(bold+cyan, line))
	fmt.Println(l.c(bold+cyan, "  oze  ") + l.c(bold, title))
	fmt.Println(l.c(bold+cyan, line))
	fmt.Println()
}

// Step prints an iteration header.
func (l *Logger) Step(i, max int, msg string) {
	ts := l.c(gray, time.Now().Format("15:04:05"))
	iter := l.c(bold+blue, fmt.Sprintf("[%d/%d]", i, max))
	fmt.Printf("\n%s %s %s\n", ts, iter, msg)
}

// Info prints an informational message.
func (l *Logger) Info(msg string) {
	fmt.Println(l.c(gray, "  ℹ ") + msg)
}

// Warn prints a warning.
func (l *Logger) Warn(msg string) {
	fmt.Println(l.c(yellow, "  ⚠ ") + msg)
}

// Error prints an error.
func (l *Logger) Error(msg string) {
	fmt.Println(l.c(red, "  ✖ ") + msg)
}

// Success prints a success line.
func (l *Logger) Success(msg string) {
	fmt.Println(l.c(green, "  ✔ "+msg))
}

// Fail prints a failure line.
func (l *Logger) Fail(msg string) {
	fmt.Println(l.c(red, "  ✖ "+msg))
}

// Box prints a labeled output block.
func (l *Logger) Box(label, content string) {
	fmt.Println(l.c(gray, "  ┌─ "+label))
	for _, line := range strings.Split(content, "\n") {
		fmt.Println(l.c(gray, "  │ ") + line)
	}
	fmt.Println(l.c(gray, "  └─"))
}

// Done prints the final success summary.
func (l *Logger) Done(feature string, iterations int) {
	line := strings.Repeat("─", 60)
	fmt.Println()
	fmt.Println(l.c(bold+green, line))
	fmt.Println(l.c(bold+green, "  ✔ Done!"))
	fmt.Printf("  Feature:    %s\n", feature)
	fmt.Printf("  Iterations: %d\n", iterations)
	fmt.Println(l.c(bold+green, line))
	fmt.Println()
}

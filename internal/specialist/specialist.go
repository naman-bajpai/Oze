// Package specialist defines role-based system prompts for Claude.
package specialist

import "fmt"

// Role represents a named specialist context.
type Role string

const (
	Frontend  Role = "frontend"
	Backend   Role = "backend"
	Mobile    Role = "mobile"
	Database  Role = "database"
	DevOps    Role = "devops"
	Security  Role = "security"
)

// prompts maps each role to a concise system prompt appended to Claude's context.
var prompts = map[Role]string{
	Frontend: "You are a frontend specialist." +
		" Focus on: React/Next.js, TypeScript, Tailwind CSS, accessibility (WCAG), responsive design, Core Web Vitals." +
		" Prefer CSS solutions over JS where possible. Keep components small and composable." +
		" Do not touch backend files unless absolutely necessary.",

	Backend: "You are a backend specialist." +
		" Focus on: REST/GraphQL API design, authentication, input validation, error handling, scalability, and security." +
		" Follow the existing framework conventions. Do not touch frontend files unless absolutely necessary.",

	Mobile: "You are a mobile specialist." +
		" Focus on: React Native / Expo best practices, platform-specific UX (iOS vs Android), performance on low-end devices," +
		" offline-first patterns, and native module integration." +
		" Do not touch backend or web files unless absolutely necessary.",

	Database: "You are a database specialist." +
		" Focus on: schema design, indexes, query optimisation, migrations, and data integrity." +
		" Prefer additive migrations. Never drop columns or tables without explicit instruction.",

	DevOps: "You are a DevOps specialist." +
		" Focus on: CI/CD pipelines, Docker, infrastructure-as-code, secrets management, and observability." +
		" Prefer idempotent, declarative config. Do not modify application source code unless absolutely necessary.",

	Security: "You are a security specialist." +
		" Focus on: OWASP Top 10, input sanitisation, auth/authz, secrets handling, dependency vulnerabilities, and secure defaults." +
		" Flag any existing issues you encounter even if not directly related to the task.",
}

// All returns every registered role name.
func All() []Role {
	roles := make([]Role, 0, len(prompts))
	for r := range prompts {
		roles = append(roles, r)
	}
	return roles
}

// Prompt returns the system prompt for the given role.
// Returns ("", false) if the role is not registered.
func Prompt(r Role) (string, bool) {
	p, ok := prompts[r]
	return p, ok
}

// Validate returns an error if r is not a known role.
func Validate(r Role) error {
	if _, ok := prompts[r]; !ok {
		return fmt.Errorf("unknown specialist %q — valid options: frontend, backend, mobile, database, devops, security", r)
	}
	return nil
}

package cli

import "testing"

func TestAppName(t *testing.T) {
	if AppName != "echo space ai" {
		t.Errorf("expected AppName to be %q, got %q", "echo space ai", AppName)
	}
}

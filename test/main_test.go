package test

import (
	"os/exec"
	"testing"
)

func TestMainCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "-C", "../hippod", ".")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to execute main: %v\nOutput:\n%s", err, string(out))
	}
}
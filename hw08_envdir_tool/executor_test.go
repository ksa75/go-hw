package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// RunCommand captures the output of running the go-envdir command.
func RunCommand(envDir, command string, args ...string) (string, error) {
	cmdArgs := append([]string{envDir, command}, args...)
	cmd := exec.Command("./go-envdir", cmdArgs...)
	cmd.Env = append(os.Environ(), "ADDED=from original env")

	// Capture output
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestRunCmd(t *testing.T) {
	// Prepare test data
	envDir := "./testdata/env"
	script := "./testdata/echo.sh"
	args := []string{"arg1=1", "arg2=2"}

	// Run the command
	result, err := RunCommand(envDir, "/bin/bash", script, args...)
	if err != nil {
		t.Fatalf("Failed to run go-envdir: %v", err)
	}

	// Define the expected output
	expected := `HELLO is (hello)
BAR is (bar)
FOO is (foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2`

	// Compare the result with the expected output
	if strings.TrimSpace(result) != expected {
		t.Errorf("Invalid output: %s\nExpected: %s", result, expected)
	}
}

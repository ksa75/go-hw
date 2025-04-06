package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// RunCommand captures the output of running the go-envdir command.
func RunCommand(envDir, command string, args ...string) (string, error) {
	cmdArgs := append([]string{envDir, command}, args...)
	fmt.Println(cmdArgs)
	cmd := exec.Command("./go-envdir", cmdArgs...)
	cmd.Env = append(os.Environ(), "ADDED=from original env")

	// Capture output
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestRunCmd(t *testing.T) {
	envDir := "./testdata/env"
	command := "/bin/bash ./testdata/echo.sh"
	args := []string{"arg1=1", "arg2=2"}

	result, err := RunCommand(envDir, command, args...)
	if err != nil {
		t.Fatalf("Failed to run go-envdir: %v", err)
	}

	// Define the expected output
	expected := `HELLO is (hello)
BAR is (bar)
FOO is (   foo
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

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	defer os.Remove("go-envdir")

	envDir := "./testdata/env"
	command := "./testdata/echo.sh"
	args := []string{"arg1=1", "arg2=2"}
	cmdArgs := append([]string{envDir, command}, args...)

	cmd := exec.Command("./go-envdir", cmdArgs...)
	cmd.Env = append(os.Environ(), "ADDED=from original env")

	result, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run go-envdir: %v", err)
	}

	// Define the expected output
	expected := `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2`

	// Compare the result with the expected output
	if strings.TrimSpace(string(result)) != expected {
		t.Errorf("Invalid output: %s\nExpected: %s", result, expected)
	}
}

func init() {
	args := []string{"build", "-o", "go-envdir"}
	cmd := exec.Command("go", args...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("building failed", string(result), err)
	}
}

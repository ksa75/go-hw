package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, value := range env {
		if value.NeedRemove {
			// If NeedRemove is true, unset the variable
			syscall.Unsetenv(key)
		} else {
			// Otherwise, set the environment variable
			syscall.Setenv(key, value.Value)
		}
	}

	// Prepare command and arguments
	command := cmd[0]
	args := cmd[1:]

	// Create a new command with the given arguments
	cmdExec := exec.Command(command, args...)

	// Redirect standard I/O to the parent process
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Stdin = os.Stdin

	// Run the command
	err := cmdExec.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			fmt.Println("Command failed with exit code:", exitError.ExitCode())
		} else {
			fmt.Println("Unexpected error:", err)
		}
	}

	// Return the exit code of the command
	return cmdExec.ProcessState.ExitCode()
}

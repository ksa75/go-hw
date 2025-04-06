package main

import (
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
		// If an error occurs, capture the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		} else {
			fmt.Printf("Error executing command: %v\n", err)
			return 1
		}
	}

	// If successful, return 0 exit code
	return 0
}

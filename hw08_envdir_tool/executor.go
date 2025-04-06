package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Prepare environment variables
	for key, value := range env {
		if value.NeedRemove {
			// If NeedRemove is true, unset the variable
			if runtime.GOOS == "windows" {
				// On Windows, use os.Unsetenv
				os.Unsetenv(key)
			} else {
				// On Linux and other OSs, use syscall.Unsetenv
				syscall.Unsetenv(key)
			}
		} else {
			// Otherwise, set the environment variable
			if runtime.GOOS == "windows" {
				// On Windows, use os.Setenv
				os.Setenv(key, value.Value)
			} else {
				// On Linux and other OSs, use syscall.Setenv
				syscall.Setenv(key, value.Value)
			}
		}
	}

	// Prepare the command
	command := cmd[0]
	args := cmd[1:]

	// If on Windows, prepend cmd.exe /C to run the command in a shell
	if runtime.GOOS == "windows" {
		command = "cmd.exe"
		args = append([]string{"/C"}, cmd...)
	}

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

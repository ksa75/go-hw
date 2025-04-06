package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Ensure correct number of arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: go-envdir <envdir> <command> [args...]")
		os.Exit(1)
	}

	// Get the directory containing the environment files and the command to run
	envDir := os.Args[1]
	cmd := os.Args[2:]
	// fmt.Println(runtime.GOOS)

	// Read environment variables from the specified directory
	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatalf("Error reading environment directory: %v", err)
	}

	// Run the command with the environment variables
	exitCode := RunCmd(cmd, env)

	// Exit with the same code as the command
	os.Exit(exitCode)
}

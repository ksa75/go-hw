package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %s, %v", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fPath := dir + "/" + file.Name()
		f, err := os.Open(fPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %s, %v", fPath, err)
		}
		defer f.Close()

		// Read the first line from the file
		scanner := bufio.NewScanner(f)
		if scanner.Scan() {
			value := strings.TrimSpace(scanner.Text())
			if value == "" {
				// If the file is empty, mark for removal
				env[file.Name()] = EnvValue{NeedRemove: true}
			} else {
				// Otherwise, store the value
				env[file.Name()] = EnvValue{Value: value}
			}
		} else if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file %s: %v", fPath, err)
		}
	}

	return env, nil
}

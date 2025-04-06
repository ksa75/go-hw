package main

import (
	"testing"
)

func TestReadDir(t *testing.T) {
	// Prepare test data
	testEnvDir := "./testdata/env"

	// Call ReadDir
	env, err := ReadDir(testEnvDir)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	// Define the expected environment variables from the test files
	expectedEnv := map[string]EnvValue{
		"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
		"BAR":   {Value: "bar", NeedRemove: false},
		"HELLO": {Value: "\"hello\"", NeedRemove: false},
		"UNSET": {NeedRemove: true}, // This should be removed as its file is empty
		"EMPTY": {NeedRemove: true}, // This should be removed as its file is empty
	}

	// Check that the environment variables match the expected values
	for key, expectedValue := range expectedEnv {
		envValue, exists := env[key]
		if !exists {
			t.Errorf("Expected environment variable %s to exist", key)
			continue
		}
		if envValue != expectedValue {
			t.Errorf("For %s, expected %v, but got %v", key, expectedValue, envValue)
		}
	}
}

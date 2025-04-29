package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	SliceInt struct {
		Val []int `validate:"min:1|max:10"`
	}

	SliceString struct {
		Val []string `validate:"regexp:[abc"`
	}

	SomeString struct {
		Val string `validate:"regexp:[abc"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid user",
			in: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "John",
				Age:    30,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "09876543210"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid user (age, email, role, phone length)",
			in: User{
				ID:     "12345678-1234-1234-1234-123456789012",
				Name:   "Alice",
				Age:    17,
				Email:  "invalid_email",
				Role:   "user",
				Phones: []string{"12345"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age"},
				{Field: "Email"},
				{Field: "Role"},
				{Field: "Phones[0]"},
			},
		},
		{
			name: "good slice of ints",
			in: SliceInt{
				Val: []int{1, 2, 3},
			},
			expectedErr: nil,
		},
		{
			name: "bad slice of ints",
			in: SliceInt{
				Val: []int{1, 2, 30},
			},
			expectedErr: ValidationErrors{
				{Field: "Val[2]"},
			},
		},

		{
			name: "slice regexp compile error",
			in: SliceString{
				Val: []string{"неважно что", "неважно как"},
			},
			expectedErr: errors.New("invalid regexp"),
		},
		{
			name: "string regexp compile error",
			in: SomeString{
				Val: "неважно что",
			},
			expectedErr: errors.New("invalid regexp"),
		},

		{
			name: "valid app",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "invalid response code",
			in: Response{
				Code: 302,
				Body: "Redirect",
			},
			expectedErr: ValidationErrors{
				{Field: "Code"},
			},
		},
		{
			name: "valid response code",
			in: Response{
				Code: 404,
				Body: "Not found",
			},
			expectedErr: nil,
		},
		{
			name:        "non-struct input",
			in:          42,
			expectedErr: errors.New("only structs are supported"),
		},
		{
			name:        "token struct without validations",
			in:          Token{},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)

			checkErrors(t, err, tt.expectedErr)
		})
	}
}

func checkErrors(t *testing.T, actualErr, expectedErr error) {
	t.Helper()

	if expectedErr == nil && actualErr != nil {
		t.Errorf("unexpected error: %v", actualErr)
		return
	}
	if expectedErr != nil && actualErr == nil {
		t.Errorf("expected error but got nil")
		return
	}
	if expectedErr == nil && actualErr == nil {
		// оба nil - всё хорошо
		return
	}

	if isValidationErrors(expectedErr) {
		checkValidationErrors(t, actualErr, expectedErr)
		return
	}

	fmt.Println()
	fmt.Println("!!!!   Program Error: ", actualErr)
	fmt.Println()
	checkProgramErrors(t, actualErr, expectedErr)
}

func isValidationErrors(err error) bool {
	var ve ValidationErrors
	return errors.As(err, &ve)
}

func checkValidationErrors(t *testing.T, actualErr, expectedErr error) {
	t.Helper()

	var expected ValidationErrors
	var actual ValidationErrors

	if !errors.As(expectedErr, &expected) {
		t.Errorf("expected error should be ValidationErrors")
		return
	}
	if !errors.As(actualErr, &actual) {
		t.Errorf("actual error should be ValidationErrors")
		return
	}

	if len(expected) != len(actual) {
		t.Errorf("expected %d validation errors, got %d", len(expected), len(actual))
		return
	}

	for idx := range expected {
		if expected[idx].Field != actual[idx].Field {
			t.Errorf("expected error in field %q, got %q", expected[idx].Field, actual[idx].Field)
		}
	}
}

func checkProgramErrors(t *testing.T, actualErr, expectedErr error) {
	t.Helper()

	if !strings.Contains(actualErr.Error(), expectedErr.Error()) {
		t.Errorf("expected error containing %q, got %q", expectedErr.Error(), actualErr.Error())
	}
}

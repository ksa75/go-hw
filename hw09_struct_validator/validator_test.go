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
				Age:    17,                // min:18
				Email:  "invalid_email",   // bad regexp
				Role:   "user",            // not in admin,stuff
				Phones: []string{"12345"}, // len!=11
			},
			expectedErr: ValidationErrors{
				{Field: "Age"},
				{Field: "Email"},
				{Field: "Role"},
				{Field: "Phones[0]"},
			},
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
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			// Сравниваем наличие ошибки
			if tt.expectedErr == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectedErr != nil && err == nil {
				t.Errorf("expected error but got nil")
			}

			// Если обе ошибки не nil — сравниваем содержимое
			if tt.expectedErr != nil && err != nil {
				switch expected := tt.expectedErr.(type) {
				case ValidationErrors:
					actual, ok := err.(ValidationErrors)
					if !ok {
						t.Errorf("expected ValidationErrors type, got: %T", err)
					} else {
						if len(expected) != len(actual) {
							t.Errorf("expected %d validation errors, got %d", len(expected), len(actual))
						} else {
							// По полям сравнение
							for idx := range expected {
								if expected[idx].Field != actual[idx].Field {
									t.Errorf("expected field error '%s', got '%s'", expected[idx].Field, actual[idx].Field)
								}
							}
						}
					}
				default:
					if !strings.Contains(err.Error(), tt.expectedErr.Error()) {
						t.Errorf("expected error containing %q, got %q", tt.expectedErr.Error(), err.Error())
					}
				}
			}
			_ = tt
		})
	}
}

package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
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

func validUser() User {
	return User{
		ID:     strings.Repeat("a", 36),
		Age:    25,
		Email:  "user@mail.com",
		Role:   "admin",
		Phones: []string{"89123456789", "89876543210"},
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          validUser(),
			expectedErr: nil,
		},
		{
			in:          func() *User { u := validUser(); return &u }(),
			expectedErr: nil,
		},
		{
			in: func() User { u := validUser(); u.Age = 10; return u }(),
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("value must be greater than %d", 18),
				},
			},
		},
		{
			in: func() User { u := validUser(); u.Email = "wrong"; return u }(),
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("regexp does not match %s", "^\\w+@\\w+\\.\\w+$"),
				},
			},
		},
		{
			in: func() User { u := validUser(); u.Role = "guest"; return u }(),
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Role",
					Err:   fmt.Errorf("%s is not in %s", "in", "admin,stuff"),
				},
			},
		},
		{
			in: func() User { u := validUser(); u.Age = 60; return u }(),
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("value must be less than %d", 50),
				},
			},
		},
		{
			in: func() User { u := validUser(); u.ID = "short"; return u }(),
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   fmt.Errorf("expected length %v, got %v", 36, 5),
				},
			},
		},
		{
			in: func() User { u := validUser(); u.Phones[1] = "123"; return u }(),
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Phones[1]",
					Err:   fmt.Errorf("expected length %v, got %v", 11, 3),
				},
			},
		},
		{
			in:          App{Version: "1.2.3"},
			expectedErr: nil,
		},
		{
			in: App{Version: "long"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   fmt.Errorf("expected length %v, got %v", 5, 4),
				},
			},
		},
		{
			in:          Response{Code: 200},
			expectedErr: nil,
		},
		{
			in: Response{Code: 201},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   fmt.Errorf("value must contain %s", "200,404,500"),
				},
			},
		},
		{
			in:          Token{Header: []byte("h"), Payload: []byte("p")},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)

			require.Equal(t, tt.expectedErr, err)
		})
	}
}

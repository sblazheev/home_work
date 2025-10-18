package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	Meta struct {
		Update string `validate:"len:19"`
		Enable bool
		Index  string `validate:"len:6|regexp:^\\d+$"`
	}
	MetaErRules struct {
		Update string `validate:"len19"`
		Enable bool
		Index  string `validate:"len:6|regexp:^\\d+$"`
	}
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   Meta     `validate:"nested"`
	}

	UserMeta struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		Meta   Meta     `validate:"nested"`
		app    App      `validate:"nested"`
	}

	UserErRules struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int         `validate:"min:|max:50"`
		Email  string      `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole    `validate:"in:admin,stuff"`
		Phones []string    `validate:"len11"`
		Meta   MetaErRules `validate:"nested"`
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

var tests = []struct {
	in          interface{}
	expectedErr error
}{
	{
		in: UserErRules{
			ID:     "4e33a976-b726-4377-a8f3-5ac93e190bfd",
			Name:   "Name",
			Age:    30,
			Email:  "test@test.ru",
			Role:   "admin",
			Phones: []string{"74951234567"},
			Meta: MetaErRules{
				Update: "2025-10-10 10:10:59",
				Index:  "111111",
			},
		},
		expectedErr: ValidatorErrors{
			ValidatorError{
				Field: "Age",
				Err:   ErrFormatRule,
			},
			ValidatorError{
				Field: "Phones",
				Err:   ErrFormatRule,
			},
			ValidatorError{
				Field: "Meta",
				Err: ValidatorErrors{
					ValidatorError{
						Field: "Update",
						Err:   ErrFormatRule,
					},
				},
			},
		},
	},
	{
		in: User{
			ID:     "111",
			Name:   "Name",
			Age:    99,
			Email:  "Email",
			Role:   "Role",
			Phones: []string{"123456"},
			meta:   Meta{},
		},
		expectedErr: ValidationErrors{
			ValidationError{
				Field: "ID",
				Err:   ErrValidationStringLength,
			},
			ValidationError{
				Field: "Age",
				Err:   ErrValidationIntNotMax,
			},
			ValidationError{
				Field: "Email",
				Err:   ErrValidationRegExpNotMatch,
			},
			ValidationError{
				Field: "Role",
				Err:   ErrValidationNotIncludes,
			},
			ValidationError{
				Field: "Phones",
				Err:   ErrValidationStringLength,
			},
		},
	},
	{
		in: UserMeta{
			ID:     "4e33a976-b726-4377-a8f3-5ac93e190bfd",
			Name:   "Name",
			Age:    30,
			Email:  "test@test.ru",
			Role:   "admin",
			Phones: []string{"74951234567"},
			Meta: Meta{
				Update: "2025-10-10 10:10:59",
				Index:  "111111",
			},
			app: App{
				Version: "1",
			},
		},
		expectedErr: nil,
	},
	{
		in: UserMeta{
			ID:     "4e33a976-b726-4377-a8f3-5ac93e190bfd",
			Name:   "Name",
			Age:    30,
			Email:  "test@test.ru",
			Role:   "admin",
			Phones: []string{"74951234567"},
			Meta: Meta{
				Update: "2025-10-10 10:10:59",
				Index:  "wwwwww",
			},
			app: App{
				Version: "1.3.4",
			},
		},
		expectedErr: ValidationErrors{
			ValidationError{
				Field: "Meta",
				Err: ValidationErrors{
					ValidationError{
						Field: "Index",
						Err:   ErrValidationRegExpNotMatch,
					},
				},
			},
		},
	},
	{
		in: User{
			ID:     "4e33a976-b726-4377-a8f3-5ac93e190bfd",
			Name:   "Name",
			Age:    30,
			Email:  "test@test.ru",
			Role:   "admin",
			Phones: []string{"74951234567"},
		},
		expectedErr: nil,
	},
	{
		in: App{
			Version: "1.34",
		},
		expectedErr: ValidationErrors{
			ValidationError{
				Field: "Version",
				Err:   ErrValidationStringLength,
			},
		},
	},
	{
		in: App{
			Version: "1.3.4",
		},
		expectedErr: nil,
	},
	{
		in: Token{
			Header:  []byte{0, 1, 2},
			Payload: nil,
		},
		expectedErr: nil,
	},
	{
		in: Token{
			Header:    []byte("application: text/json"),
			Payload:   []byte("{\"id\":1}"),
			Signature: []byte("signature"),
		},
		expectedErr: nil,
	},
	{
		in: Response{
			Code: 204,
			Body: "",
		},
		expectedErr: ValidationErrors{
			ValidationError{
				Field: "Code",
				Err:   ErrValidationNotIncludes,
			},
		},
	},
	{
		in: Response{
			Code: 200,
			Body: "OK",
		},
		expectedErr: nil,
	},
}

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

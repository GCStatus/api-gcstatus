package tests

import (
	"errors"
	"gcstatus/internal/utils"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := map[string]struct {
		key       string
		token     string
		expectErr bool
	}{
		"valid encryption and decryption": {
			key:       "myverystrongpasswordo32bitlength",
			token:     "testtoken123",
			expectErr: false,
		},
		"short key": {
			key:       "shortkey",
			token:     "testtoken123",
			expectErr: true,
		},
		"empty token": {
			key:       "myverystrongpasswordo32bitlength",
			token:     "",
			expectErr: false,
		},
		"empty key": {
			key:       "",
			token:     "testtoken123",
			expectErr: true,
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			encryptedToken, err := utils.Encrypt(test.token, test.key)
			if (err != nil) != test.expectErr {
				t.Errorf("Encrypt(%q, %q) returned error %v; expected error: %v", test.token, test.key, err, test.expectErr)
			}

			if !test.expectErr {
				decryptedToken, err := utils.Decrypt(encryptedToken, test.key)
				if err != nil {
					t.Errorf("Failed to decrypt token: %v", err)
				}

				if decryptedToken != test.token {
					t.Errorf("Expected %q, got %q after decryption", test.token, decryptedToken)
				}
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := map[string]struct {
		input  string
		result bool
	}{
		"valid password": {
			input:  "ValidPass123!",
			result: true,
		},
		"short string": {
			input:  "short1!",
			result: false,
		},
		"all as lower case": {
			input:  "alllowercase123!",
			result: false,
		},
		"all as upper case": {
			input:  "ALLUPPERCASE123!",
			result: false,
		},
		"without digit or symbols": {
			input:  "NoDigitsOrSymbols",
			result: false,
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got, expected := utils.ValidatePassword(test.input), test.result; got != expected {
				t.Fatalf("validatePassword(%q) returned %v; expected %v", test.input, got, expected)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := map[string]struct {
		password  string
		expectErr bool
	}{
		"valid password": {
			password:  "StrongPass123!",
			expectErr: false,
		},
		"empty password": {
			password:  "",
			expectErr: true,
		},
		"special characters": {
			password:  "P@ssw0rd!",
			expectErr: false,
		},
		"whitespace": {
			password:  "  ",
			expectErr: true,
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			hashedPassword, err := utils.HashPassword(test.password)

			if (err != nil) != test.expectErr {
				t.Errorf("HashPassword(%q) returned error: %v; expected error: %v", test.password, err, test.expectErr)
			}

			if !test.expectErr {
				err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(test.password))
				if err != nil {
					t.Errorf("Password does not match the hash: %v", err)
				}
			}
		})
	}
}

func TestGetFirstAndLastName(t *testing.T) {
	tests := map[string]struct {
		fullName  string
		firstName string
		lastName  string
	}{
		"normal names": {
			fullName:  "John Doe",
			firstName: "John",
			lastName:  "Doe",
		},
		"single name": {
			fullName:  "Alice",
			firstName: "Alice",
			lastName:  "",
		},
		"multiple names": {
			fullName:  "Alice Bob Charlie",
			firstName: "Alice",
			lastName:  "Charlie",
		},
		"empty string": {
			fullName:  "",
			firstName: "",
			lastName:  "",
		},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			firstName, lastName := utils.GetFirstAndLastName(test.fullName)
			if firstName != test.firstName || lastName != test.lastName {
				t.Errorf("Expected (%s, %s), got (%s, %s)", test.firstName, test.lastName, firstName, lastName)
			}
		})
	}
}

func IsHashEqualsValueTest(t *testing.T) {
	base := "admin1234"
	baseInvalid := "abcdefg"

	hashed, err := utils.HashPassword(base)
	if err != nil {
		t.Fatalf("failed to hash the base password: %s", err.Error())
	}

	tests := map[string]struct {
		hash        string
		value       string
		expectEqual bool
	}{
		"is equals": {
			hash:        hashed,
			value:       base,
			expectEqual: true,
		},
		"not equals": {
			hash:        hashed,
			value:       baseInvalid,
			expectEqual: false,
		},
		"error case with invalid hash": {
			hash:        "invalid-hash",
			value:       base,
			expectEqual: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			equal := utils.IsHashEqualsValue(tc.hash, tc.value)

			assert.NoError(t, err, "did not expect an error")
			assert.Equal(t, tc.expectEqual, equal, "expectEqual check failed")
		})
	}
}

func TestFormatValidationError(t *testing.T) {
	type SampleStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	validate := validator.New()

	tests := map[string]struct {
		input        error
		expected     []string
		prepareInput func() error
	}{
		"as input validation error": {
			prepareInput: func() error {
				sample := SampleStruct{}
				return validate.Struct(sample)
			},
			expected: []string{
				"Name is required and cannot be empty.",
				"Email is required and cannot be empty.",
			},
		},
		"as a generic error": {
			input: errors.New("some generic error"),
			expected: []string{
				"some generic error",
			},
		},
		"as nil input": {
			input:    nil,
			expected: []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var err error

			if tc.prepareInput != nil {
				err = tc.prepareInput()
			} else {
				err = tc.input
			}

			formattedErrors := utils.FormatValidationError(err)

			assert.Equal(t, tc.expected, formattedErrors)
		})
	}
}

func TestFormatTimestamp(t *testing.T) {
	fixedTime := time.Now()

	tests := map[string]struct {
		date time.Time
	}{
		"valid date": {
			date: fixedTime,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			timestamp := utils.FormatTimestamp(tc.date)

			assert.Equal(t, timestamp, fixedTime.Format("2006-01-02T15:04:05"))
		})
	}
}

func TestNormalizeWhitespace(t *testing.T) {
	tests := map[string]struct {
		toNormalize  string
		expectReturn string
	}{
		"can normalize string": {
			toNormalize: `
				Hello,
				World.
				This is a test!
			`,
			expectReturn: "Hello, World. This is a test!",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			normalized := utils.NormalizeWhitespace(tc.toNormalize)

			assert.Equal(t, normalized, tc.expectReturn)
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"simple string": {
			input:    "Hello World",
			expected: "hello-world",
		},
		"string with special chars": {
			input:    "Hello World!",
			expected: "hello-world",
		},
		"string with numbers": {
			input:    "123 Hello World",
			expected: "123-hello-world",
		},
		"multiple spaces": {
			input:    "Hello    World",
			expected: "hello----world",
		},
		"leading and trailing spaces": {
			input:    "  Hello World  ",
			expected: "hello-world",
		},
		"all caps": {
			input:    "HELLO WORLD",
			expected: "hello-world",
		},
		"empty string": {
			input:    "",
			expected: "",
		},
		"only special characters": {
			input:    "@#$%",
			expected: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := utils.Slugify(tc.input)
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

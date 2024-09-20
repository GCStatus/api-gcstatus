package tests

import (
	"gcstatus/pkg/utils"
	"testing"

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

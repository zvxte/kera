package model

import (
	"strings"
	"testing"
)

func TestIsUsernameValid(t *testing.T) {
	validUsernames := []string{
		"AbcD", "1234", "abcd", "ab_cd", "_1234", "AB__CD", "ABCDabcd12345678",
	}
	for _, username := range validUsernames {
		if !IsUsernameValid(username) {
			t.Errorf("username should be valid: %q", username)
		}
	}

	invalidUsernames := []string{
		"A", "Ab", "Abc", "123", "_abc_", "____", "ABCDabcd123456789", "abcd.", "abcdπ", "ab cd",
	}
	for _, username := range invalidUsernames {
		if IsUsernameValid(username) {
			t.Errorf("username should not be valid: %q", username)
		}
	}
}

func TestIsDisplayNameValid(t *testing.T) {
	validDisplayNames := []string{
		"AbcD", "1234", "AB CD", "ABCDabcd12345678", "!@#$",
	}
	for _, displayName := range validDisplayNames {
		if !IsDisplayNameValid(displayName) {
			t.Errorf("display name should be valid: %q", displayName)
		}
	}

	invalidDisplayNames := []string{
		"A", "Abc", "!@#", "ABCDabcd123456789", "πππ", "    ", "ab  c", " abcd",
	}
	for _, displayName := range invalidDisplayNames {
		if IsDisplayNameValid(displayName) {
			t.Errorf("display name should not be valid: %q", displayName)
		}
	}
}

func TestIsPlainPasswordValid(t *testing.T) {
	validPlainPasswords := []string{
		"12345678", strings.Repeat("1", plainPasswordMaxChars),
	}
	for _, plainPassword := range validPlainPasswords {
		if !IsPlainPasswordValid(plainPassword) {
			t.Errorf("plain password should be valid: %q", plainPassword)
		}
	}

	invalidPlainPasswords := []string{
		"1", "1234", "1234567", strings.Repeat("1", plainPasswordMaxChars+1),
	}
	for _, plainPassword := range invalidPlainPasswords {
		if IsPlainPasswordValid(plainPassword) {
			t.Errorf("plain password should not be valid: %q", plainPassword)
		}
	}
}

func TestIsHashedPasswordValid(t *testing.T) {
	validHashedPasswords := []string{
		"$argon2id$v=19$m=19,t=2,p=1$MTIzNDU2Nzg$4zj0kIPamvrQt6pWuGm2rg",
		"$argon2id$v=19$m=19,t=2,p=1$MTIzNDU2Nzg$re2ihnZD/LvtAnf/gaol5+VqV/ENlwmkdudQUA+qpNs",
		strings.Repeat("1", hashedPasswordMaxChars),
	}
	for _, hashedPassword := range validHashedPasswords {
		if !IsHashedPasswordValid(hashedPassword) {
			t.Errorf("hashed password should be valid: %q", hashedPassword)
		}
	}

	invalidHashedPasswords := []string{
		"1", "1234567", strings.Repeat("1", hashedPasswordMaxChars+1),
	}
	for _, hashedPassword := range invalidHashedPasswords {
		if IsHashedPasswordValid(hashedPassword) {
			t.Errorf("hashed password should not be valid: %q", hashedPassword)
		}
	}
}

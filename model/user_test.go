package model

import (
	"strings"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	validUsernames := []string{
		"AbcD", "1234", "abcd", "ab_cd", "_1234", "AB__CD", "ABCDabcd12345678",
	}
	for _, username := range validUsernames {
		if err := validateUsername(username); err != nil {
			t.Errorf("username %q should be valid: %v", username, err)
		}
	}

	invalidUsernames := []string{
		"A", "Ab", "Abc", "123", "_abc_", "____", "ABCDabcd123456789", "abcd.", "abcdπ", "ab cd",
	}
	for _, username := range invalidUsernames {
		if err := validateUsername(username); err == nil {
			t.Errorf("username %q should not be valid: %v", username, err)
		}
	}
}

func TestValidateDisplayName(t *testing.T) {
	validDisplayNames := []string{
		"AbcD", "1234", "AB CD", "ABCDabcd12345678", "!@#$",
	}
	for _, displayName := range validDisplayNames {
		if err := validateDisplayName(displayName); err != nil {
			t.Errorf("display name %q should be valid: %v", displayName, err)
		}
	}

	invalidDisplayNames := []string{
		"A", "Abc", "!@#", "ABCDabcd123456789", "πππ", "    ", "ab  c", " abcd",
	}
	for _, displayName := range invalidDisplayNames {
		if err := validateDisplayName(displayName); err == nil {
			t.Errorf("display name %q should not be valid: %v", displayName, err)
		}
	}
}

func TestValidatePlainPassword(t *testing.T) {
	validPlainPasswords := []string{
		"12345678", strings.Repeat("1", plainPasswordMaxChars),
	}
	for _, plainPassword := range validPlainPasswords {
		if err := validatePlainPassword(plainPassword); err != nil {
			t.Errorf("plain password %q should be valid: %v", plainPassword, err)
		}
	}

	invalidPlainPasswords := []string{
		"1", "1234", "1234567", strings.Repeat("1", plainPasswordMaxChars+1),
	}
	for _, plainPassword := range invalidPlainPasswords {
		if err := validatePlainPassword(plainPassword); err == nil {
			t.Errorf("plain password %q should not be valid: %v", plainPassword, err)
		}
	}
}

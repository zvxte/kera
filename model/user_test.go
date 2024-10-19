package model

import "testing"

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
	for _, username := range validDisplayNames {
		if !IsDisplayNameValid(username) {
			t.Errorf("display name should be valid: %q", username)
		}
	}

	invalidDisplayNames := []string{
		"A", "Abc", "!@#", "ABCDabcd123456789", "πππ", "    ", "ab  c", " abcd",
	}
	for _, username := range invalidDisplayNames {
		if IsDisplayNameValid(username) {
			t.Errorf("display name should not be valid: %q", username)
		}
	}
}

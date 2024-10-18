package model

import (
	"fmt"
	"testing"
)

func TestIsUsernameValid(t *testing.T) {
	validUsernames := []string{
		"AbcD", "1234", "abcd", "ab_c", "_123", "AB__", "ABCDabcd12345678",
	}
	for _, username := range validUsernames {
		if !IsUsernameValid(username) {
			t.Error(fmt.Sprintf("Username should be valid: %q", username))
		}
	}

	invalidUsernames := []string{
		"A", "Ab", "Abc", "123", "____", "ABCDabcd123456789", "abcd.", "abcdπ", "ab cd",
	}
	for _, username := range invalidUsernames {
		if IsUsernameValid(username) {
			t.Error(fmt.Sprintf("Username should not be valid: %q", username))
		}
	}
}

package session

import "testing"

func TestValidateID(t *testing.T) {
	tests := []struct {
		name      string
		sessionID string
		expected  bool
	}{
		{
			"Valid",
			"2IBEnhqxKuWE5dhXF4IaQJrBrSygJnRq",
			true,
		},
		{
			"Invalid: too short",
			"2IBEnhqxKuWE5dhXF4IaQJrBrSygJnR",
			false,
		},
		{
			"Invalid: too long",
			"2IBEnhqxKuWE5dhXF4IaQJrBrSygJnRqX",
			false,
		},
		{
			"Invalid: invalid characters",
			"!IBEnhqxKuWE5dhXF4IaQJrBrSygJnR@",
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ValidateID(test.sessionID)
			if got != test.expected {
				t.Errorf(
					"ValidateID(%q), got=%v, expected=%v",
					test.sessionID, got, test.expected,
				)
			}
		})
	}
}

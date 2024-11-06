package model

import "testing"

func TestNewSessionID(t *testing.T) {
	s, err := NewSessionID()
	if err != nil {
		t.Error(err, s)
	}
}

func TestValidateSessionID(t *testing.T) {
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
			got := ValidateSessionID(test.sessionID)
			if got != test.expected {
				t.Errorf(
					"ValidateSessionID(%q), got=%v, expected=%v",
					test.sessionID, got, test.expected,
				)
			}
		})
	}
}

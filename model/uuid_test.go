package model

import "testing"

func TestNewUUIDv7(t *testing.T) {
	if _, err := NewUUIDv7(); err != nil {
		t.Error(err)
	}
}

func TestParseUUID(t *testing.T) {
	validValue := "0192a14e-0605-726d-89ed-6b39dae2ea91"
	if _, err := ParseUUID(validValue); err != nil {
		t.Error(err)
	}

	invalidValue := ""
	if _, err := ParseUUID(invalidValue); err == nil {
		t.Error(err)
	}
}

func TestUUID(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		validValue := "0192a14e-0605-726d-89ed-6b39dae2ea91"
		id, _ := ParseUUID(validValue)
		if result := id.String(); result != validValue {
			t.Errorf("expected UUID string %q, but got %q", validValue, result)
		}
	})
}

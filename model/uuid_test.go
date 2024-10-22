package model

import "testing"

func TestNewUUIDv7(t *testing.T) {
	if _, err := NewUUIDv7(); err != nil {
		t.Error(err)
	}
}

func TestParseUUID(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		shouldErr bool
	}{
		{"Valid", "0192a14e-0605-726d-89ed-6b39dae2ea91", false},
		{"Invalid", "0192a14e-", true},
		{"Invalid", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseUUID(test.value)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"ParseUUID(%q), error=%v, shouldErr=%v",
					test.value, err, test.shouldErr,
				)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		id       UUID
		shouldBe string
	}{
		{
			"Valid all zeros",
			UUID{},
			"00000000-0000-0000-0000-000000000000",
		},
		{
			"Valid random",
			UUID{
				0x55, 0x0e, 0x84, 0x00,
				0xe2, 0x9b, 0x41, 0xd4,
				0xa7, 0x16, 0x44, 0x66,
				0x55, 0x44, 0x00, 0x00,
			},
			"550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.id.String()
			if result != test.shouldBe {
				t.Errorf(
					"UUID.String(), got=%q, expected=%q",
					result, test.shouldBe,
				)
			}
		})
	}
}

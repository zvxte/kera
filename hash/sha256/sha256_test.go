package sha256

import (
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		shouldBe [32]byte
	}{
		{
			"Valid",
			"sessionID",
			[32]byte{176, 83, 141, 118, 44, 177, 23, 75, 116, 170, 183, 4, 185, 39, 75, 0, 48, 194, 98, 68, 125, 55, 148, 124, 252, 61, 93, 99, 226, 201, 36, 192},
		},
		{
			"Valid: short",
			"a",
			[32]byte{202, 151, 129, 18, 202, 27, 189, 202, 250, 194, 49, 179, 154, 35, 220, 77, 167, 134, 239, 248, 20, 124, 78, 114, 185, 128, 119, 133, 175, 238, 72, 187},
		},
		{
			"Valid: long",
			strings.Repeat("a", 256),
			[32]byte{2, 215, 22, 13, 119, 225, 140, 100, 71, 190, 128, 194, 227, 85, 199, 237, 67, 136, 84, 82, 113, 112, 44, 80, 37, 59, 9, 20, 198, 92, 229, 254},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash := Hash(test.input)
			if hash != test.shouldBe {
				t.Errorf(
					"Hash(%q), got=%v, expected=%v",
					test.input, hash, test.shouldBe,
				)
			}
		})
	}
}

func TestVerifyHash(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		hashedInput [32]byte
		shouldBe    bool
	}{
		{
			"Valid",
			"sessionID",
			[32]byte{176, 83, 141, 118, 44, 177, 23, 75, 116, 170, 183, 4, 185, 39, 75, 0, 48, 194, 98, 68, 125, 55, 148, 124, 252, 61, 93, 99, 226, 201, 36, 192},
			true,
		},
		{
			"Valid",
			"aaaaaaaa",
			[32]byte{31, 60, 228, 4, 21, 162, 8, 31, 163, 238, 231, 95, 195, 159, 255, 142, 86, 194, 34, 112, 209, 169, 120, 167, 36, 155, 89, 45, 206, 189, 32, 180},
			true,
		},
		{
			"Invalid",
			"AAAAaaaa",
			[32]byte{0, 1, 2, 3, 4},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := VerifyHash(test.input, test.hashedInput)
			if isValid != test.shouldBe {
				t.Errorf(
					"VerifyHash(%q, %q), got=%v, expected=%v",
					test.input, test.hashedInput, isValid, test.shouldBe,
				)
			}
		})
	}
}

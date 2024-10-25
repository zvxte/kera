package argon2id

import (
	"bytes"
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"Valid", "password", false},
		{"Valid short", "a", false},
		{"Valid long", strings.Repeat("a", 256), false},
		{"Valid empty", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Hash(test.input, DefaultParams)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"Hash(%q, %q), error=%v, shouldErr=%v",
					test.input, DefaultParams, err, test.shouldErr,
				)
			}
		})
	}
}

func TestVerifyHash(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		hashedInput string
		shouldErr   bool
		shouldBe    bool
	}{
		{
			"Valid",
			"password",
			"$argon2id$v=19$m=19456,t=2,p=1$YWFhYUFBQUFhYWFhQUFBQQ$KdIUCTl6NPY+m4WM+pHJW0fWIQMLQV5LCZE7zYvqryU",
			false,
			true,
		},
		{
			"Valid short",
			"a",
			"$argon2id$v=19$m=19456,t=4,p=2$QUFBQWFhYWFBQUFBYWFhYQ$IPFpgW+DUFhT3f007pRzz2HUV5wPAV0VzANxSuo81g4",
			false,
			true,
		},
		{
			"Valid long",
			strings.Repeat("a", 256),
			"$argon2id$v=19$m=9,t=4,p=1$MTIzNDU2Nzg$HjQikypjR1bPBW7IkAIwi3Khxu4HLjRBwl1KBRDf/w4",
			false,
			true,
		},
		{
			"Invalid",
			"password",
			"$argon2id$v=19$m=19,t=2,p=1$",
			true,
			false,
		},
		{
			"Invalid",
			"password",
			"$argon2id$v=19$m=19,t=2,p=1$MTIzNDU2Nzg$",
			true,
			false,
		},
		{
			"Invalid empty",
			"password",
			"",
			true,
			false,
		},
		{
			"Invalid variant",
			"password",
			"$argon2i$v=19$m=19,t=2,p=1$MTIzNDU2Nzg$5vI4/d3YW0ADXglN8ziuIJoqS/dj3wNFLOcc394xvRk",
			true,
			false,
		},
		{
			"Invalid version",
			"password",
			"$argon2i$v=18$m=19,t=2,p=1$MTIzNDU2Nzg$5vI4/d3YW0ADXglN8ziuIJoqS/dj3wNFLOcc394xvRk",
			true,
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid, err := VerifyHash(test.input, test.hashedInput)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"VerifyHash(%q, %q), error=%v, shouldErr=%v",
					test.input, test.hashedInput, err, test.shouldErr,
				)
			}
			if isValid != test.shouldBe {
				t.Errorf(
					"VerifyHash(%q, %q), got=%v, expected=%v",
					test.input, test.hashedInput, isValid, test.shouldBe,
				)
			}
		})
	}
}

func TestDecodeHash(t *testing.T) {
	tests := []struct {
		name           string
		hashedInput    string
		expectedParams *Params
		expectedSalt   []byte
		expectedKey    []byte
		shouldErr      bool
	}{
		{
			"Valid",
			"$argon2id$v=19$m=19456,t=2,p=1$YWFhYUFBQUFhYWFhQUFBQQ$sxzziMCgbNOhfrgUXet7cS5rE2gq2pLe5hUaXLC966I",
			DefaultParams,
			[]byte("aaaaAAAAaaaaAAAA"),
			[]byte{179, 28, 243, 136, 192, 160, 108, 211, 161, 126, 184, 20, 93, 235, 123, 113, 46, 107, 19, 104, 42, 218, 146, 222, 230, 21, 26, 92, 176, 189, 235, 162},
			false,
		},
		{
			"Valid",
			"$argon2id$v=19$m=19456,t=4,p=2$YmJiYkJCQkJiYmJiQkJCQg$z43l5zj14v2qQnPFLWzS+Ci6HFoA3IbezmvPtVnCTAk",
			&Params{
				memory:      19 * 1024,
				iterations:  4,
				parallelism: 2,
				keyLength:   32,
				saltLength:  16,
			},
			[]byte("bbbbBBBBbbbbBBBB"),
			[]byte{207, 141, 229, 231, 56, 245, 226, 253, 170, 66, 115, 197, 45, 108, 210, 248, 40, 186, 28, 90, 0, 220, 134, 222, 206, 107, 207, 181, 89, 194, 76, 9},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			params, salt, key, err := decodeHash(test.hashedInput)
			if (err != nil) != test.shouldErr {
				t.Errorf(
					"decodeHash(%q), error=%v, shouldErr=%v",
					test.hashedInput, err, test.shouldErr,
				)
			}
			if *params != *test.expectedParams {
				t.Errorf(
					"decodeHash(%q), params=%v, expectedParams=%v",
					test.hashedInput, params, test.expectedParams,
				)
			}
			if !bytes.Equal(salt, test.expectedSalt) {
				t.Errorf(
					"decodeHash(%q), salt=%v, expectedSalt=%v",
					test.hashedInput, salt, test.expectedSalt,
				)
			}
			if !bytes.Equal(key, test.expectedKey) {
				t.Errorf(
					"decodeHash(%q), key=%v, expectedKey=%v",
					test.hashedInput, key, test.expectedKey,
				)
			}
		})
	}
}

func TestGenerateSalt(t *testing.T) {
	_, err := generateSalt(DefaultParams.saltLength)
	if err != nil {
		t.Error(err)
	}
}

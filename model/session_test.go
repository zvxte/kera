package model

import "testing"

func TestNewSessionID(t *testing.T) {
	s, err := NewSessionID()
	if err != nil {
		t.Error(err, s)
	}
}

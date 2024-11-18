package session

import "testing"

func TestNewID(t *testing.T) {
	s, err := NewID()
	if err != nil {
		t.Error(err, s)
	}
}

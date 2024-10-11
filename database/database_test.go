package database

import (
	"testing"
)

func TestNewSqlDatabase(t *testing.T) {
	invalid := ""
	ErrorExpectedMessage := "expected an error, but got nil"

	_, err := NewSqlDatabase(invalid, invalid)
	if err == nil {
		t.Error(ErrorExpectedMessage)
	}

	_, err = NewSqlDatabase(Postgres, invalid)
	if err == nil {
		t.Error(ErrorExpectedMessage)
	}
}

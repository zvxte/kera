package database

import (
	"testing"
)

func TestNewSqlDatabase(t *testing.T) {
	invalid := ""
	ErrorExpectedMessage := "expected an error, but got nil"

	_, err := NewSqlDatabase(InvalidDriverName, invalid)
	if err == nil {
		t.Error(ErrorExpectedMessage)
	}

	_, err = NewSqlDatabase(PostgresDriverName, invalid)
	if err == nil {
		t.Error(ErrorExpectedMessage)
	}
}

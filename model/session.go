package model

import "time"

type Session struct {
	ID             string
	CreationDate   time.Time
	ExpirationDate time.Time
}

package model

import "time"

type User struct {
	ID            UUID
	Username      string
	UsernameLower string
	DisplayName   string
	Password      string
	Location      time.Location
	CreationDate  time.Time
}

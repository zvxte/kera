package store

import (
	"github.com/zvxte/kera/model"
)

type UserStore interface {
	CreateUser(user model.User) error
	GetUserByUsername(username string) (model.User, error)
	GetUserByUUID(uuid model.UUID) (model.User, error)
}

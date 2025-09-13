package repo

import (
	"core/models"
)

type UserRepo interface {
	// create defualt user
	CreateDefaultUser(hasPassword string) error
	// create a new user
	CreateUser(user *models.User) error
}

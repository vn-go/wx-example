package core

import (
	"core/models"
)

type userRepo interface {
	// create defualt user
	CreateDefaultUser(hasPassword string) error
	// create a new user
	CreateUser(user *models.User) error
}

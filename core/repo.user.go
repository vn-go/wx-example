package core

import (
	"context"
	"core/models"
)

type userRepo interface {
	// create defualt user
	CreateDefaultUser(ctx context.Context, hasPassword string) error
	// create a new user
	CreateUser(ctx context.Context, user *models.User) error
}

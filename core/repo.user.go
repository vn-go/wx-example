package core

import (
	"context"
	"core/models"

	"github.com/vn-go/dx"
)

type userRepo interface {
	// create defualt user
	CreateDefaultUser(db *dx.DB, ctx context.Context, hasPassword string) error
	// create a new user
	CreateUser(db *dx.DB, ctx context.Context, user *models.User) error
	GetUserByUserId(db *dx.DB, ctx context.Context, userId string) (*models.User, error)
	GetUserByName(db *dx.DB, ctx context.Context, username string) (*models.User, error)
	DeleteUserByUserId(db *dx.DB, ctx context.Context, userId string) error
}

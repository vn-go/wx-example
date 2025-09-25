package core

import (
	"context"
	"sync"

	"core/models"

	"github.com/vn-go/dx"
)

type userRepoSql struct {
}

func (repo *userRepoSql) CreateUser(db *dx.DB, ctx context.Context, user *models.User) error {
	return db.WithContext(ctx).Insert(user)
}
func (repo *userRepoSql) GetUserByUserId(db *dx.DB, ctx context.Context, userId string) (*models.User, error) {
	ret := &models.User{}
	if err := db.WithContext(ctx).First(ret, "userId=?", userId); err != nil {
		if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
			if dbErr.ErrorType == dx.Errors.NOTFOUND {
				return nil, nil
			}
		}
		return nil, err
	}
	return ret, nil
}
func (repo *userRepoSql) GetUserByName(db *dx.DB, ctx context.Context, username string) (*models.User, error) {
	ret := &models.User{}
	if err := db.WithContext(ctx).First(ret, "username=?", username); err != nil {
		return nil, err
	}
	return ret, nil
}

var CreateDefaultUserOnce sync.Once

func (repo *userRepoSql) CreateDefaultUser(db *dx.DB, ctx context.Context, hashPassword string) error {
	var err error
	CreateDefaultUserOnce.Do(func() {
		var user *models.User
		user, err = dx.NewThenSetDefaultValues(func() (*models.User, error) {
			return &models.User{
				Username:     "root",
				HashPassword: hashPassword,
				//CreatedOn:    time.Now().UTC(),
			}, nil
		})
		if err != nil {
			return
		}

		err = repo.CreateUser(db, ctx, user)
		if err != nil {

			if dxErr := dx.Errors.IsDbError(err); dxErr != nil {
				if dxErr.ErrorType != dx.Errors.DUPLICATE {
					err = dxErr
					return
				} else {
					err = nil

				}
			} else {

				return
			}
		}

	})

	return err
}
func (repo *userRepoSql) DeleteUserByUserId(db *dx.DB, ctx context.Context, userId string) error {
	return db.Delete(&models.User{}, "userId=?", userId).Error
}
func newUserRepoSql() userRepo {
	return &userRepoSql{}
}

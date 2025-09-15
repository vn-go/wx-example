package core

import (
	"context"
	"core/internal"
	"core/models"

	"github.com/vn-go/dx"
)

type userRepoSql struct {
	db      *dx.DB
	context context.Context
}

func newUserRepoSql(db *dx.DB, ctx context.Context) userRepo {
	return &userRepoSql{
		db:      db,
		context: ctx,
	}
}
func (repo *userRepoSql) CreateUser(user *models.User) error {
	return repo.db.WithContext(repo.context).Insert(user)
}

func (repo *userRepoSql) CreateDefaultUser(hashPassword string) error {
	_, err := internal.OnceCall[userRepoSql]("CreateDefaultUser", func() (*models.User, error) {
		user, err := dx.NewThenSetDefaultValues(func() (*models.User, error) {
			return &models.User{
				Username:     "root",
				HashPassword: hashPassword,
				//CreatedOn:    time.Now().UTC(),
			}, nil
		})
		if err != nil {
			return nil, err
		}

		err = repo.CreateUser(user)
		if err != nil {

			if dxErr := dx.Errors.IsDbError(err); dxErr != nil {
				if dxErr.ErrorType != dx.Errors.DUPLICATE {
					return nil, dxErr

				}
			} else {

				return nil, err
			}
		}
		return user, nil
	})

	return err
}

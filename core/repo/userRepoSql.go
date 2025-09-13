package repo

import (
	"context"
	"core/internal"
	"core/models"

	"github.com/vn-go/dx"
)

type UserRepoSql struct {
	db      *dx.DB
	context context.Context
}

func NewUserRepoSql(db *dx.DB, context context.Context) UserRepo {
	return &UserRepoSql{
		db:      db,
		context: context,
	}
}
func (repo *UserRepoSql) CreateUser(user *models.User) error {
	return repo.db.WithContext(repo.context).Insert(user)
}

func (repo *UserRepoSql) CreateDefaultUser(hashPassword string) error {
	_, err := internal.OnceCall[UserRepoSql]("CreateDefaultUser", func() (*models.User, error) {
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

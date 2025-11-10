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
	// ret := &models.User{}
	// if err := db.WithContext(ctx).First(ret, "username=?", username); err != nil {
	// 	return nil, err
	// }
	// return ret, nil
	return dx.QueryItem[models.User](db, "user(),where(username=?)", username)
}

type initCreateDefaultUser struct {
	err  error
	once sync.Once
}

var initCreateDefaultUserCache sync.Map

func (repo *userRepoSql) CreateDefaultUser(db *dx.DB, ctx context.Context, username, hashPassword string) error {

	a, _ := initCreateDefaultUserCache.LoadOrStore(db.DbName+"@"+db.DriverName, &initCreateDefaultUser{})
	i := a.(*initCreateDefaultUser)
	i.once.Do(func() {
		var user *models.User
		var err error
		user, err = dx.NewThenSetDefaultValues(func() (*models.User, error) {
			return &models.User{
				Username:     username,
				HashPassword: hashPassword,
				IsSysAdmin:   true,
				//CreatedOn:    time.Now().UTC(),
			}, nil
		})
		user.IsSysAdmin = true

		if err != nil {
			i.err = err
			return
		}

		err = repo.CreateUser(db, ctx, user)
		if err != nil {

			if dxErr := dx.Errors.IsDbError(err); dxErr != nil {
				if dxErr.ErrorType != dx.Errors.DUPLICATE {
					err = dxErr
					return
				} else {
					return

				}
			} else {
				i.err = err
				return
			}

		}

	})
	if i.err != nil {
		initCreateDefaultUserCache.Delete(db.DbName + "@" + db.DriverName)
	}
	return i.err
}
func (repo *userRepoSql) DeleteUserByUserId(db *dx.DB, ctx context.Context, userId string) error {
	return db.Delete(&models.User{}, "userId=?", userId).Error
}
func newUserRepoSql() userRepo {
	return &userRepoSql{}
}

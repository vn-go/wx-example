package repo

import (
	"cmp/wxapi/models"

	
)

type UserRepo interface {
	CreateUser(user *models.User) error
}

type UserRepoSql struct {
	db *xdb.TenantDB
}

func NewUserRepoSql(db *xdb.TenantDB) UserRepo {
	return &UserRepoSql{
		db: db,
	}
}
func (userRepo *UserRepoSql) CreateUser(user *models.User) error {
	return userRepo.db.Insert(user)
}

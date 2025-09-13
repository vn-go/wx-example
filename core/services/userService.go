package services

import (
	"context"
	"core/internal"
	"core/models"
	"core/repo"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vn-go/dx"
)

type UserService interface {
	CreateUser(user *models.User) error
}

type UserServiceSql struct {
	userRepo repo.UserRepo
	pwdSvc   PasswordService
}

var userServiceSqlCreateUserDefault sync.Once

func (userSvc *UserServiceSql) CreateUser(user *models.User) error {
	var err error

	userServiceSqlCreateUserDefault.Do(func() {
		hashPass, err := userSvc.pwdSvc.HashPassword("root@123456")
		if err != nil {
			return
		}
		err = userSvc.userRepo.CreateDefaultUser(hashPass)
	})
	if err != nil {
		return err
	}
	hashPass, err := internal.OnceCall[UserServiceSql]("CreateUser/HashPassword/Root", func() (string, error) {
		return userSvc.pwdSvc.HashPassword("root@root")

	})
	if err != nil {
		return err
	}

	userSvc.userRepo.CreateDefaultUser(hashPass)
	user.UserId = uuid.NewString()
	user.CreatedOn = time.Now().UTC()
	hashPass, err = userSvc.pwdSvc.HashPassword(fmt.Sprintf("%s@%s", strings.ToLower(user.Username), user.HashPassword))
	if err != nil {
		return err
	}
	user.HashPassword = hashPass
	return userSvc.userRepo.CreateUser(user)
}
func NewUserServiceSql(db *dx.DB, ctx context.Context) UserService {
	return &UserServiceSql{
		userRepo: repo.NewUserRepoSql(db, ctx),
		pwdSvc:   NewBcryptPasswordService(),
	}
}

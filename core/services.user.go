package core

import (
	"context"
	"core/models"

	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vn-go/bx"
)

type userService interface {
	CreateUser(tenant string, ctx context.Context, user *models.User) error
	GetUserByUserId(tenant string, ctx context.Context, userId string) (*models.User, error)
	SayHello() int
}

type userServiceSql struct {
	userRepo userRepo
	pwdSvc   passwordService
	cache    cacheService
	ctx      context.Context
	tenant   tenantService
}

var userServiceSqlCreateUserDefault sync.Once

func (userSvc *userServiceSql) CreateUser(tenant string, ctx context.Context, user *models.User) error {
	var err error
	tenantDb, err := userSvc.tenant.GetTenant(tenant)
	if err != nil {
		return err
	}
	userServiceSqlCreateUserDefault.Do(func() {
		hashPass, err := userSvc.pwdSvc.HashPassword("root", "123456")
		if err != nil {
			return
		}
		err = userSvc.userRepo.CreateDefaultUser(tenantDb, ctx, hashPass)
	})
	if err != nil {
		return err
	}

	hashPass, err := bx.OnceCall[userServiceSql]("CreateUser/HashPassword/Root", func() (string, error) {
		return userSvc.pwdSvc.HashPassword("root", "123456")

	})
	if err != nil {
		return err
	}

	userSvc.userRepo.CreateDefaultUser(tenantDb, ctx, hashPass)
	user.UserId = uuid.NewString()
	user.CreatedOn = time.Now().UTC()
	hashPass, err = userSvc.pwdSvc.HashPassword(user.Username, user.HashPassword)
	if err != nil {
		return err
	}
	user.HashPassword = hashPass
	if err := userSvc.userRepo.CreateUser(tenantDb, ctx, user); err != nil {
		return err
	}

	if err := userSvc.cache.Set(userSvc, ctx, user.UserId+"@user", user); err != nil {
		return err
	}
	user.HashPassword = ""
	return nil
}
func (userSvc *userServiceSql) GetUserByUserId(tenant string, ctx context.Context, userId string) (*models.User, error) {
	var retUser *models.User
	if err := userSvc.cache.Get(userSvc, ctx, userId+"@user", retUser); err == nil {
		return retUser, nil
	}
	tenantDb, err := userSvc.tenant.GetTenant(tenant)
	if err != nil {
		return nil, err
	}
	retUser, err = userSvc.userRepo.GetUserByUserId(tenantDb, ctx, userId)
	if err != nil {
		return nil, err
	}
	if retUser != nil {
		if err := userSvc.cache.Set(userSvc, ctx, retUser.UserId+"@user", retUser); err != nil {
			return nil, err
		}
	}
	return retUser, nil
}
func (userSvc *userServiceSql) SayHello() int {
	return 1 + 1
}
func newUserServiceSql(
	tenant tenantService,
	cache cacheService,
	userRepo userRepo,
	passwordSvc passwordService,
) userService {
	return &userServiceSql{
		userRepo: userRepo,
		pwdSvc:   passwordSvc,
		cache:    cache,
		tenant:   tenant,
	}
}

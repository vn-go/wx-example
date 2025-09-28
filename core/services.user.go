package core

import (
	"context"
	"core/models"
	"strings"

	"sync"
	"time"

	"github.com/google/uuid"
)

type userService interface {
	CreateUser(tenant string, ctx context.Context, user *models.User) error
	GetUserByUserId(tenant string, ctx context.Context, userId string) (*models.User, error)
	DeleteUserByUserId(tenant string, ctx context.Context, userId string) error
	AddUser(tenanet string, ctx context.Context, user *models.User) error
}

type userServiceSql struct {
	userRepo userRepo
	pwdSvc   passwordService
	cache    cacheService

	tenant tenantService
}

var userServiceSqlCreateUserDefault sync.Once

func (userSvc *userServiceSql) CreateUser(tenant string, ctx context.Context, user *models.User) error {
	var err error
	var hashPass string
	tenantDb, err := userSvc.tenant.GetTenant(tenant)
	if err != nil {
		return err
	}

	// userServiceSqlCreateUserDefault.Do(func() {
	// 	hashPass, err = userSvc.pwdSvc.HashPassword("root", "123456")
	// 	if err != nil {
	// 		return
	// 	}
	// 	err = userSvc.userRepo.CreateDefaultUser(tenantDb, ctx, hashPass)
	// })
	// if err != nil {
	// 	return err
	// }

	// hashPass, err = bx.OnceCall[userServiceSql]("CreateUser/HashPassword/Root", func() (string, error) {
	// 	return userSvc.pwdSvc.HashPassword("root", "123456")

	// })
	// if err != nil {
	// 	return err
	// }

	// userSvc.userRepo.CreateDefaultUser(tenantDb, ctx, hashPass)
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

	if err := userSvc.cache.AddObject(ctx, tenant, user.UserId, user, 4); err != nil {
		return err
	}
	user.HashPassword = ""
	return nil
}
func (userSvc *userServiceSql) GetUserByUserId(tenant string, ctx context.Context, userId string) (*models.User, error) {
	var retUser *models.User
	if err := userSvc.cache.GetObject(ctx, tenant, userId, retUser); err == nil {
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
		if err := userSvc.cache.AddObject(ctx, tenant, userId, retUser, 4); err != nil {
			return nil, err
		}
	}
	return retUser, nil
}
func (userSvc *userServiceSql) DeleteUserByUserId(tenant string, ctx context.Context, userId string) error {
	tenantDb, err := userSvc.tenant.GetTenant(tenant)
	if err != nil {
		return err
	}
	err = userSvc.userRepo.DeleteUserByUserId(tenantDb, ctx, userId)
	if err != nil {
		return err
	}
	user := &models.User{}
	if err := userSvc.cache.GetObject(ctx, tenant, userId, user); err != nil {
		return err
	}

	//delete auth cache item
	cacheItem := &OAuthResponseCacheItem{}
	if err := userSvc.cache.DeleteObject(ctx, tenant, strings.ToLower(user.Username), cacheItem); err == nil {
		return err
	}
	//delete cache compare pass
	if err := userSvc.cache.DeleteObject(ctx, tenant, strings.ToLower(user.Username), &ComparePasswordCacheItem{}); err != nil {
		return err
	}
	// then delete cache usser
	if err := userSvc.cache.DeleteObject(ctx, tenant, userId, &models.User{}); err != nil {
		return err
	}

	return nil
}
func (userSvc *userServiceSql) AddUser(tenanet string, ctx context.Context, user *models.User) error {
	panic("implete me")
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

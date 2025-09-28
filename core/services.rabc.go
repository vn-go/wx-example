package core

import (
	"context"
	"core/models"
	"time"

	"github.com/google/uuid"
	"github.com/vn-go/dx"
)

type rabcService interface {
	NewRole(ctx context.Context, creator *UserClaims, role *models.Role) (*models.Role, error)
	NewUser(ctx context.Context, creator *UserClaims, user *models.User) (*models.User, error)
	//GetDb(user *UserClaims) (*dx.DB, error)
}

type rabcServiceImpl struct {
	tanentSvc tenantService
	cache     cacheService
}

func (s *rabcServiceImpl) NewRole(ctx context.Context, user *UserClaims, role *models.Role) (*models.Role, error) {
	var db *dx.DB
	var err error
	if db, err = s.tanentSvc.GetTenant(user.Tenant); err != nil {
		return nil, err
	}
	role.CreatedBy = user.Username
	role.CreatedOn = time.Now().UTC()
	role.RoleId = uuid.NewString()
	role.Description = dx.Ptr(dx.IsNull((*role).Description, ""))

	err = db.WithContext(ctx).Insert(role)
	return role, err

}
func (s *rabcServiceImpl) NewUser(ctx context.Context, creator *UserClaims, user *models.User) (*models.User, error) {
	user.CreatedBy = creator.Username
	var db *dx.DB
	var err error
	if db, err = s.tanentSvc.GetTenant(creator.Tenant); err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Insert(user)
	return user, err
}

//	func (s *rabcServiceImpl) GetDb(user *UserClaims) (*dx.DB, error) {
//		if s.cfg.Tenant.IsMulti {
//			return s.tanentSvc.GetTenant(user.Tenant)
//		}
//		return s.tanentSvc.GetTenant(user.Tenant)
//	}
func NewRabcService(
	tanentSvc tenantService,
	cache cacheService,
) rabcService {
	return &rabcServiceImpl{
		tanentSvc: tanentSvc,
		cache:     cache,
	}
}

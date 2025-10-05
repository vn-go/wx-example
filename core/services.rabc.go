package core

import (
	"context"
	"core/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vn-go/dx"
)

type Pager struct {
	Index   int      `json:"index" check:"range:[0:)"`
	Size    int      `json:"size" check:"range:[10:)"`
	OrderBy []string `json:"orderBy"`
}
type AccountInfo struct {
	//models.Role
	Username string `json:"username"`
	UserId   string `json:"userId"`

	RoleId   *string `json:"roleId"`
	RoleCode *string `json:"roleCode"`
	RoleName *string `json:"roleName"`
}
type rabcService interface {
	NewRole(ctx context.Context, creator *UserClaims, role *models.Role) (*models.Role, error)
	GetRoleByRoleId(ctx context.Context, user *UserClaims, roleId string) (*models.Role, error)
	NewUser(ctx context.Context, creator *UserClaims, user *models.User) (*models.User, error)
	GetListOfRoles(ctx context.Context, user *UserClaims, pager Pager) ([]models.Role, error)
	//GetListOfRolesSQL(ctx context.Context, user *UserClaims, pager Pager) (string, error)
	GetListOfAccounts(ctx context.Context, user *UserClaims, pager Pager) ([]AccountInfo, error)
	ChangeUserPassword(ctx context.Context, user *UserClaims, username string, newPass string) error
	//GetDb(user *UserClaims) (*dx.DB, error)
}

type rabcServiceImpl struct {
	tanentSvc tenantService
	pwdSvc    passwordService
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
func (s *rabcServiceImpl) GetRoleByRoleId(ctx context.Context, user *UserClaims, roleId string) (role *models.Role, err error) {

	var db *dx.DB

	if db, err = s.tanentSvc.GetTenant(user.Tenant); err != nil {
		return nil, err
	}
	role = &models.Role{}
	err = db.WithContext(ctx).First(role, "roleId=?", roleId)
	return role, err
}

// func (s *rabcServiceImpl) GetListOfRolesSQL(ctx context.Context, user *UserClaims, pager Pager) (sql string, err error) {
// 	var db *dx.DB

// 	if db, err = s.tanentSvc.GetTenant(user.Tenant); err != nil {
// 		return "", err
// 	}
// 	// roles = []models.Role{}

// 	qr := db.WithContext(ctx).Limit(uint64(pager.Size)).Offset(uint64(pager.Index * pager.Size))
// 	if len(pager.OrderBy) > 0 {
// 		qr = qr.Order(strings.Join(pager.OrderBy, ","))
// 	}
// 	sqlP, err := qr..ToSql(db, reflect.TypeOf(models.Role{}))
// 	return sqlP.Sql, err

// }
func (s *rabcServiceImpl) GetListOfRoles(ctx context.Context, user *UserClaims, pager Pager) (roles []models.Role, err error) {
	var db *dx.DB

	if db, err = s.tanentSvc.GetTenant(user.Tenant); err != nil {
		return nil, err
	}
	roles = []models.Role{}

	qr := db.WithContext(ctx).Limit(uint64(pager.Size)).Offset(uint64(pager.Index * pager.Size))
	if len(pager.OrderBy) > 0 {
		qr = qr.Order(strings.Join(pager.OrderBy, ","))
	}
	err = qr.Find(&roles)
	return roles, err

}
func (s *rabcServiceImpl) GetListOfAccounts(ctx context.Context, user *UserClaims, pager Pager) (accs []AccountInfo, err error) {
	var db *dx.DB

	if db, err = s.tanentSvc.GetTenant(user.Tenant); err != nil {
		return nil, err
	}
	//dx.Options.ShowSql = true
	qr := db.WithContext(ctx).From(&models.User{}).Joins("u left join role r on u.roleId=r.id").Select(
		`u.userId UserId,
		u.username Username,
		r.code RoleCode,
		r.name RoleName,
		r.roleId RoleId`,
	)
	qr.Limit(uint64(pager.Size))
	qr.Offset(uint64(pager.Size * pager.Index))
	accs = []AccountInfo{}
	err = qr.Find(&accs)
	return accs, err

}
func (s *rabcServiceImpl) ChangeUserPassword(ctx context.Context, user *UserClaims, username string, newPass string) (err error) {
	var db *dx.DB

	if db, err = s.tanentSvc.GetTenant(user.Tenant); err != nil {
		return err
	}
	updateUser := &models.User{}
	if err = db.First(updateUser, "username=?", username); err != nil {
		return err
	}
	hashPassword, err := s.pwdSvc.HashPassword(username, newPass)
	if err != nil {
		return err
	}
	updateUser.HashPassword = hashPassword
	r := db.Update(updateUser)
	return r.Error

}
func NewRabcService(
	tanentSvc tenantService,
	cache cacheService,
	pwdSvc passwordService,

) rabcService {
	return &rabcServiceImpl{
		tanentSvc: tanentSvc,
		cache:     cache,
		pwdSvc:    pwdSvc,
	}
}

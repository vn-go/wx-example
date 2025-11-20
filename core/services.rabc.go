package core

import (
	"context"
	"core/models"
	"crypto/sha256"
	"encoding/hex"
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
type DataSouceViewInfo struct {
	ViewPath string `json:"viewPath"`
	Dsl      string
}

type rabcService interface {
	NewRole(ctx context.Context, creator *UserClaims, role *models.Role) (*models.Role, error)
	GetRoleByRoleId(ctx context.Context, user *UserClaims, roleId string) (*models.Role, error)
	NewUser(ctx context.Context, creator *UserClaims, user *models.User) (*models.User, error)
	GetListOfRoles(ctx context.Context, user *UserClaims, pager Pager) (any, error)
	//GetListOfRolesSQL(ctx context.Context, user *UserClaims, pager Pager) (string, error)
	GetListOfAccounts(ctx context.Context, user *UserClaims, pager Pager) (any, error)
	GetAccountById(ctx context.Context, user *UserClaims, userId string) (any, error)
	//GetDb(user *UserClaims) (*dx.DB, error)
	ResgisterDataSouceView(ctx context.Context, user *UserClaims, Info DataSouceViewInfo)
	ResgisterView(ctx context.Context, Tenant, viewPath, apiPath, createdBy string) error
}

type rabcServiceImpl struct {
	tanentSvc tenantService
	pwdSvc    passwordService
	cache     cacheService
}

/*
this function get account info by userid
*/
func (r rabcServiceImpl) GetAccountById(ctx context.Context, user *UserClaims, userId string) (any, error) {
	var db *dx.DB
	var err error
	if db, err = r.tanentSvc.GetTenant(user.Tenant); err != nil {
		return nil, err
	}

	return dx.QueryItem[models.User](db, "user(),where(userId=?)", userId)
}
func (r rabcServiceImpl) ResgisterView(ctx context.Context, Tenant, viewPath, apiPath, createdBy string) error {
	db, err := r.tanentSvc.GetTenant(Tenant)
	if err != nil {
		return err
	}
	//make sha256 ViewPath
	sum := sha256.Sum256([]byte(viewPath)) // trả về [32]byte
	ViewId := hex.EncodeToString(sum[:])
	view, err := dx.NewDTO[models.UIView]()
	if err != nil {
		return err
	}
	view.ViewId = ViewId
	view.ViewPath = viewPath
	view.CreatedBy = createdBy
	view.CreatedOn = time.Now().UTC()

	err = db.InsertWithContext(ctx, view)
	if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
		if dbErr.ErrorType != dx.Errors.DUPLICATE {
			return err
		} else {

			view, err = dx.NewQuery[models.UIView](
				"uiView(id)", //<-- get id column only
			).Filter("viewPath=?", viewPath).ToItem(db)
			// err = db.DslFirstRow(view, "uiview(id),where(viewId=?)", ViewId)
			if err != nil {
				return err
			}
		}
	}
	api, err := dx.NewDTO[models.Api]()
	if err != nil {
		return err
	}
	api.ApiPath = apiPath
	api.CreatedOn = time.Now().UTC()
	err = db.InsertWithContext(ctx, api)
	if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
		if dbErr.ErrorType != dx.Errors.DUPLICATE {
			return err
		} else {
			api, err = dx.NewQuery[models.Api]("api(id)").Filter("apiPath=?", apiPath).ToItemWithContext(ctx, db)
			if err != nil {
				return err
			}
		}
	}
	viewApi, err := dx.NewDTO[models.UIViewApi]()
	if err != nil {
		return err
	}
	viewApi.ViewId = view.Id
	viewApi.ApiId = api.Id
	viewApi.CreatedBy = createdBy
	viewApi.CreatedOn = time.Now().UTC()
	err = db.InsertWithContext(ctx, viewApi)
	if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
		if dbErr.ErrorType != dx.Errors.DUPLICATE {
			return err
		}
	}
	return nil

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
	hashPass, err := s.pwdSvc.HashPassword(user.Username, user.HashPassword)
	if err != nil {
		return nil, err
	}
	user.HashPassword = hashPass
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

func (s *rabcServiceImpl) GetListOfRoles(ctx context.Context, user *UserClaims, pager Pager) (roles any, err error) {
	var db *dx.DB

	if db, err = s.tanentSvc.GetTenant(user.Tenant); err != nil {
		return nil, err
	}

	query := db.QueryModel(models.Role{}).LeftJoin(models.User{}, "role.id=user.roleId").Select(
		"role(roleId,code,name,description,createdOn),user(count(id) NumOfUsers)",
	)
	if len(pager.OrderBy) == 0 {
		query = query.SortDesc("role.createdOn")
	} else {
		query = query.Sort(strings.Join(pager.OrderBy, ","))
	}
	return query.ToArray()

}
func (s *rabcServiceImpl) GetListOfAccounts(ctx context.Context, user *UserClaims, pager Pager) (items any, err error) {
	var db *dx.DB

	if db, err = s.tanentSvc.GetTenant(user.Tenant); err != nil {
		return nil, err
	}

	query := db.QueryModel(models.User{}).LeftJoin(models.Role{}, "user.roleId=role.id").Select(
		`user(username,userId,email,createdOn,isActive),
			role(code RoleCode,name RoleName,description),
			take(?)`,
		pager.Size,
	)
	if len(pager.OrderBy) == 0 {
		query = query.SortDesc("user.createdOn")
	} else {
		query = query.Sort(strings.Join(pager.OrderBy, ","))
	}
	return query.ToArray()

}

func (s *rabcServiceImpl) ResgisterDataSouceView(ctx context.Context, user *UserClaims, Info DataSouceViewInfo) {
	panic("implement me")
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

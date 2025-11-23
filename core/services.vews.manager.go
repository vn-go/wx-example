package core

import (
	"context"
	"core/models"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/vn-go/dx"
)

type viewService struct {
	tenantDb tenantService
}

func (v *viewService) DeleteApi(context context.Context, user *UserClaims, path string) (any, error) {
	db, err := v.getDbByUser(user)
	if err != nil {
		return nil, err
	}
	apiItem := &models.Api{}
	err = db.First(apiItem, "apiPath=?", path)
	if err != nil {
		return nil, err
	}
	ret := db.Delete(&models.UIViewApi{}, "apiId=?", apiItem.Id)
	if ret.Error != nil {
		return nil, ret.Error
	}
	ret = db.Delete(&models.Api{}, "apiPath=?", path)
	if ret.Error != nil {
		return nil, ret.Error
	} else {
		return ret.RowsAffected, nil
	}
}

func (v *viewService) getDbByUser(user *UserClaims) (*dx.DB, error) {
	return v.tenantDb.GetTenant(user.Tenant)
}
func (v *viewService) GetListApiOfView(context context.Context, user *UserClaims, path string) (any, error) {

	db, err := v.getDbByUser(user)
	if err != nil {
		return nil, err
	}
	return db.DslToArray(`
							from(
							UIViewApi ui,
							Api a,
							UIView v ,
							ui.apiId=a.id,
							ui.viewId=v.id),
							Api(),
							where(v.viewPath=? and a.apiPath like ?)
						`, path, "%/%")
}

type initRegistApiList struct {
	once sync.Once
	err  error
}

var initRegistApiListCache sync.Map

func (v *viewService) RegistApiList(db *dx.DB, username, path string, list []string) error {
	key := fmt.Sprintf("%s;%s;%s", username, path, strings.Join(list, ";"))
	a, _ := initRegistApiListCache.LoadOrStore(key, &initRegistApiList{})
	i := a.(*initRegistApiList)
	i.once.Do(func() {
		i.err = v.RegistApiListNoCache(db, username, path, list)
	})
	return i.err

}
func (v *viewService) RegistApiListNoCache(db *dx.DB, username, path string, list []string) error {
	viewPathItem, err := dx.NewDTO[models.UIView]()
	if err != nil { // properly, models.UIView not register yet
		return err
	}
	viewPathItem.ViewPath = path
	viewPathItem.CreatedBy = username
	viewPathItem.CreatedOn = time.Now().UTC()
	err = db.Insert(viewPathItem)
	dbError := dx.Errors.IsDbError(err)
	if err != nil && dbError != nil && dbError.ErrorType == dx.Errors.DUPLICATE {
		errGetRow := db.DslFirstRow(viewPathItem, "uiView(),where(viewPath=?)", path)
		if errGetRow != nil {
			return errGetRow
		}
	}
	for _, api := range list {
		apiItem, err := dx.NewDTO[models.Api]()
		if err != nil {
			return err
		}
		apiItem.ApiPath = api
		apiItem.CreatedBy = username
		apiItem.CreatedOn = time.Now().UTC()
		err = db.Insert(apiItem)
		dbError := dx.Errors.IsDbError(err)
		if err != nil && dbError != nil && dbError.ErrorType == dx.Errors.DUPLICATE {
			err = db.DslFirstRow(apiItem, "api(),where(apiPath=?)", api) // get Id Only
			if err != nil {
				return err
			}
		}

		apiView, err := dx.NewDTO[models.UIViewApi]()
		if err != nil {
			return err
		}
		apiView.ApiId = apiItem.Id
		apiView.ViewId = viewPathItem.Id
		apiView.CreatedBy = username
		apiView.CreatedOn = time.Now().UTC()
		err = db.Insert(apiView)
		dbError = dx.Errors.IsDbError(err)
		if err != nil && dbError != nil && dbError.ErrorType == dx.Errors.DUPLICATE {
			if dbError.ErrorType == dx.Errors.DUPLICATE {
				continue // is existing continue
			}
		}
	}
	return nil
}
func (v *viewService) ApiDiscovery(context context.Context, user *UserClaims, list []string) (any, error) {

	db, err := v.tenantDb.GetTenant(user.Tenant)
	if err != nil {
		return nil, err
	}
	if user.IsUpperUser {
		err := v.RegistApiList(db, user.Username, user.ViewPath, list)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (v *viewService) GetListOfViews(context context.Context, data *UserClaims) (any, error) {
	db, err := v.tenantDb.GetTenant(data.Tenant)
	if err != nil {
		return nil, err
	}
	// a:=models.UIView{}
	// b:=models.RoleUiView{}
	// a.Id=b.ViewId
	return db.DslToArray(`
		u(viewId,viewPath,title,description,createdOn,createdBy),
		rv(count(id) TotalRoles), // get total role
		u(count(ur.id) TotalUsers), // get total user
		from(
				uiView u, 
				roleUiView rv,
				role r, 
				user ur,
				left(u.id=rv.viewId),
				left(rv.roleid=r.id),
				left(r.id=ur.roleId)
			),
		where(u.viewPath!=? and u.viewPath!=? and u.viewPath!=? and u.viewPath!=?),
		sort(u.createdOn desc)`,
		"app", "sysadmin", "undefined", "main/uiform")
}

func newViewService(tenantDb tenantService) *viewService {
	return &viewService{tenantDb: tenantDb}
}

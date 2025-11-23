package controllers

import (
	"core"
	"time"

	"core/models"

	"github.com/vn-go/dx"
	"github.com/vn-go/wx"
)

/*
Account controller
*/
type Accounts struct {
	BaseAuthController // this controller require authenticate for all  APIs
}

func (acc *Accounts) Me(h wx.Handler) any {
	return acc.Authenticate.Data
}
func (acc *Accounts) RoleCreate(h wx.Handler, role *struct {
	Code        string  `json:"code" check:"range:[3:8]"`
	Name        string  `json:"name" check:"range:[5:50]"`
	Description *string `json:"description" check:"range:[:300]"`
}) (any, error) {
	//panic(fmt.Errorf("Not implement %T", acc))
	roleData, err := dx.NewDTO[models.Role]()
	if err != nil {
		return nil, err
	}
	roleData.Code = role.Code
	roleData.Name = role.Name
	roleData.Description = role.Description
	ret, err := core.Services.RABCSvc.NewRole(h().Req.Context(), acc.Authenticate.Data, roleData)

	if err != nil {

		return nil, wx.Errors.NewHttpError(wx.ErrInternalServerError, core.Errors.Create(acc, "RoleCreate", err)) //wx.Errors.NewHttpError(wx.ErrConflict, err)
	}

	return ret, err
}
func (acc *Accounts) UserCreate(h wx.Handler, user *struct {
	Username     string  `json:"username" check:"range:[5:20]"`
	Password     string  `json:"password" check:"range:[5:20]"`
	IsSupperUser bool    `json:"isSupperUser"`
	RoleId       *string `json:"roleId" check:"range:[36:36]"`
}) (any, error) {
	userData, err := dx.NewDTO[models.User]()
	if err != nil {
		return nil, err
	}
	userData.Username = user.Username
	userData.HashPassword = user.Password
	userData.IsSysAdmin = user.IsSupperUser
	if user.RoleId != nil {
		role, err := core.Services.RABCSvc.GetRoleByRoleId(h().Req.Context(), acc.Authenticate.Data, *user.RoleId)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, wx.Errors.NewHttpError(wx.ErrNotFound, struct {
				Message string   `json:"message"`
				Fields  []string `json:"fields"`
			}{
				Message: "Role was not found",
				Fields:  []string{"RoleId"},
			})
		}
		*user.RoleId = role.RoleId
	}
	userData, err = core.Services.RABCSvc.NewUser(h().Req.Context(), acc.Authenticate.Data, userData)
	return userData, err
}
func (acc *Accounts) GetListOfRoles(h wx.Handler, pager core.Pager) (any, error) {

	return core.Services.RABCSvc.GetListOfRoles(h().Req.Context(), acc.Authenticate.Data, pager)
}

// func (acc *Accounts) GetListOfRolesSQL(h wx.Handler, pager core.Pager) (any, error) {

//		return core.Services.RABCSvc.GetListOfRolesSQL(h().Req.Context(), acc.Authenticate.Data, pager)
//	}
type PagerInfo struct {
	First uint64 `json:"first"`
	Last  uint64 `json:"last"`
}
type getListOfAccountsResult struct {
	Items        any    `json:"itens"`
	TotalRecords uint64 `json:"totalRecords"`
}
type UserEdit core.EditClaims[models.User, struct {
	UserId     string
	Username   string
	CreatedOn  time.Time
	CreatedBy  string
	ModifiedOn *time.Time
	ModifiedBy *string
}]

func (acc *Accounts) UpdateById(h wx.Handler,
	data UserEdit) (any, error) {
	err := core.Services.DataSignSvc.SignData(h().Req.Context(), acc.Data, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (acc *Accounts) GetEdit(h wx.Handler, data struct {
	UserId string `json:"userId"`
}) (*UserEdit, error) {

	ret, err := core.Services.RABCSvc.GetAccountById(h().Req.Context(), acc.Authenticate.Data, data.UserId)
	if err != nil {
		return nil, err
	}
	retData := &UserEdit{
		Data: *ret,
	}
	err = core.Services.DataSignSvc.SignData(h().Req.Context(), acc.Data, retData)
	if err != nil {
		return nil, err
	}

	return retData, err
}
func (acc *Accounts) GetListOfAccounts(h wx.Handler, pager PagerInfo) (any, error) {
	if pager.Last == 0 {
		return make([]any, 0), nil
	}
	pageSize := pager.Last - pager.First
	pageIndex := int(pager.First / pageSize)
	ret, err := core.Services.RABCSvc.GetListOfAccounts(h().Req.Context(), acc.Authenticate.Data, core.Pager{
		Index: pageIndex,
		Size:  int(pageSize),
	})
	return ret, err
}
func (acc *Accounts) ChangeUserPassword(h wx.Handler, data struct {
	Username    string `json:"username" check:"range(3:50)"`
	NewPassword string `json:"newPassword" check:"range(3:50)"`
}) error {
	return core.Services.AuthSvc.ChangePassword(h().Req.Context(), acc.Authenticate.Data, data.NewPassword)
}
func (acc *Accounts) ChangePassword(h struct {
	wx.Handler `route:"me/change-password"`
},
	data struct {
		NewPassword string `json:"newPassword" check:"range(3:50)"`
	},
) (any, error) {
	err := core.Services.AuthSvc.ChangePassword(h.Handler().Req.Context(), acc.Authenticate.Data, data.NewPassword)
	if err != nil {
		return nil, wx.Errors.NewHttpError(wx.ErrInternalServerError, core.Errors.Create(acc, "ChangePassword", err)) //wx.Errors.NewHttpError(wx.ErrConflict, err)
	}
	return nil, nil
}
func (acc *Accounts) GetMenu(h struct {
	wx.Handler `route:"me/get-menu"`
}, data []core.MenuItem) (any, error) {
	var err error
	if acc.Authenticate.Data.IsUpperUser {
		err = core.Services.SysSvc.SyncMenu(h.Handler().Req.Context(), acc.Authenticate.Data, data)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

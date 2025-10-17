package controllers

import (
	"core"

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
func (acc *Accounts) GetListOfAccounts(h wx.Handler, pager core.Pager) (any, error) {

	return core.Services.RABCSvc.GetListOfAccounts(h().Req.Context(), acc.Authenticate.Data, pager)
}
func (acc *Accounts) ChangeUserPassword(h wx.Handler, data struct {
	Username    string `json:"username" check:"range(3:50)"`
	NewPassword string `json:"newPassword" check:"range(3:50)"`
}) error {
	return core.Services.RABCSvc.ChangeUserPassword(h().Req.Context(), acc.Authenticate.Data, data.Username, data.NewPassword)
}

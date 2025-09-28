package controllers

import (
	"core"

	"core/models"

	"github.com/vn-go/dx"
	"github.com/vn-go/wx"
)

type Accounts struct {
	BaseAuthController
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

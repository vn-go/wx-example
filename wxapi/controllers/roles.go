package controllers

import (
	"core"

	"github.com/vn-go/wx"
)

type Roles struct {
	BaseAuthController
}

func (r *Roles) GetList(h wx.Handler, pager core.Pager) (any, error) {
	return core.Services.RABCSvc.GetListOfRoles(h().Req.Context(), r.Authenticate.Data, pager)
}
func (r *Roles) GetItem(h wx.Handler, roleId string) (any, error) {
	ret, err := core.Services.RABCSvc.GetRoleById(h().Req.Context(), r.Authenticate.Data, roleId)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
func (r *Roles) UpdateItem(h wx.Handler, data core.RoleEdit) (any, error) {
	ret, err := core.Services.RABCSvc.UpdateRole(h().Req.Context(), r.Authenticate.Data, &data)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

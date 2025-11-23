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
	return core.Services.RABCSvc.GetRoleById(h().Req.Context(), r.Authenticate.Data, roleId)
}

package controllers

import (
	"core"

	"github.com/vn-go/wx"
)

type Roles struct {
}

func GetListOfRoles() {
	wx.SetAuth(&Roles{}, func(ctx *wx.HttpContext[core.UserClaims]) error {
		req := ctx.Req
		authorization := req.Header["Authorization"]

		if len(authorization) == 0 {
			return wx.Errors.NewUnauthorizedError()
		}
		//s := "xknKGvzDI-sZwlXwUo2_GVsY6ce94AC3I6qNnxnOtq655tOgbcRbRnK0fs_tEb-6yz-EtjBCC0qdeDg0xu6uiw"
		user, tenant, err := core.Services.AuthSvc.Verify(req.Context(), authorization[0])
		if err != nil || user == nil {
			return wx.Errors.NewUnauthorizedError()
		}

		ctx.Identifier = core.UserClaims{
			Username:    user.Username,
			UserId:      user.UserId,
			ClaimId:     user.Id,
			Tenant:      tenant,
			IsUpperUser: user.IsSysAdmin,
		}
		return nil
	})
	wx.HandlerPost("GetListOfRoles", func(acc *Roles, context *wx.HttpContext[core.UserClaims], pager core.Pager) (any, error) {

		return core.Services.RABCSvc.GetListOfRoles(context.Req.Context(), &context.Identifier, pager)
	})
}

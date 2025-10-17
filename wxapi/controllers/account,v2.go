package controllers

import (
	"core"
	"core/models"

	"github.com/vn-go/bx"
	"github.com/vn-go/wx"
)

type Roles struct {
}

var rolesCache []models.Role

func GetListOfRoles() {
	wx.SetAuth(&Roles{}, func(ctx *wx.HttpContext[core.UserClaims]) error {
		req := ctx.Req
		authorization := req.Header["Authorization"]

		if len(authorization) == 0 {
			return wx.Errors.NewUnauthorizedError()
		}
		//s := "xknKGvzDI-sZwlXwUo2_GVsY6ce94AC3I6qNnxnOtq655tOgbcRbRnK0fs_tEb-6yz-EtjBCC0qdeDg0xu6uiw"
		u, err := bx.OnceCall[Roles](authorization[0], func() (*core.UserClaims, error) {
			user, tenant, err := core.Services.AuthSvc.Verify(req.Context(), authorization[0])
			if err != nil || user == nil {
				return nil, wx.Errors.NewUnauthorizedError()
			}

			return &core.UserClaims{
				Username:    user.Username,
				UserId:      user.UserId,
				ClaimId:     user.Id,
				Tenant:      tenant,
				IsUpperUser: user.IsSysAdmin,
			}, nil

		})
		if err != nil {
			return err
		}
		ctx.Identifier = *u

		return nil
	})
	wx.HandlerPost("GetListOfRoles", func(acc *Roles, context *wx.HttpContext[core.UserClaims], pager core.Pager) (any, error) {
		// if rolesCache == nil {
		// 	roles, err := core.Services.RABCSvc.GetListOfRoles(context.Req.Context(), &context.Identifier, pager)
		// 	rolesCache = roles
		// 	return roles, err
		// } else {
		// 	return rolesCache, nil
		// }
		return core.Services.RABCSvc.GetListOfRoles(context.Req.Context(), &context.Identifier, pager)

	})
}

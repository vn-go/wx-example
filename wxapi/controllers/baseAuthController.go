package controllers

import (
	"core"

	"github.com/vn-go/wx"
)

type BaseAuthController struct {
	wx.Authenticate[core.UserClaims]
}

func (base *BaseAuthController) New() error {
	base.Authenticate.Verify(func(ctx wx.Handler) (*core.UserClaims, error) {
		req := ctx().Req
		authorization := req.Header["Authorization"]

		if len(authorization) == 0 {
			return nil, wx.Errors.NewUnauthorizedError()
		}
		//s := "xknKGvzDI-sZwlXwUo2_GVsY6ce94AC3I6qNnxnOtq655tOgbcRbRnK0fs_tEb-6yz-EtjBCC0qdeDg0xu6uiw"
		user, tenant, err := core.Services.AuthSvc.Verify(req.Context(), authorization[0])
		if err != nil || user == nil {
			return nil, wx.Errors.NewUnauthorizedError()
		}
		return &core.UserClaims{
			Username: user.Username,
			UserId:   user.UserId,
			ClaimId:  user.Id,
			Tenant:   tenant,
		}, nil
	})
	return nil
}

type ResposeData[T any, TErr any] struct {
	Data *T
	Err  *TErr
}

func NewResposeData[T any, TErr any](data T, err TErr) *ResposeData[T, TErr] {
	return &ResposeData[T, TErr]{
		Data: &data,
		Err:  &err,
	}
}

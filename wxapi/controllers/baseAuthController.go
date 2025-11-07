package controllers

import (
	"core"

	"github.com/vn-go/wx"
)

// BaseAuthController 基础认证控制器
// All other controllers should inherit from this controller if need authentication
type BaseAuthController struct {
	wx.Authenticate[core.UserClaims]
}

// Controller constructor (auto called by framework)
func (base *BaseAuthController) New() error {
	//setup authentication middleware
	base.Authenticate.Verify(func(ctx wx.Handler) (*core.UserClaims, error) {
		ctxHandler := ctx()
		req := ctxHandler.Req

		authorization := req.Header["Authorization"]

		if len(authorization) == 0 {
			return nil, wx.Errors.NewUnauthorizedError()
		}

		user, tenant, err := core.Services.AuthSvc.Verify(req.Context(), authorization[0])
		if err != nil || user == nil {
			return nil, wx.Errors.NewUnauthorizedError()
		}
		viewPath := ""
		if req.Header["View-Path"] != nil {
			viewPath = req.Header["View-Path"][0]
		}
		ok, err := core.Services.AuthSvc.Authorize(req.Context(), tenant, user, viewPath, ctxHandler.ApiPath)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, wx.Errors.NewForbidenError("Access denied")
		}
		return &core.UserClaims{
			Username:    user.Username,
			UserId:      user.UserId,
			ClaimId:     user.Id,
			Tenant:      tenant,
			IsUpperUser: user.IsSysAdmin,
			ViewPath:    viewPath,
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

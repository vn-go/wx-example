package controllers

import (
	"context"
	"core"
	"strings"

	"github.com/vn-go/wx"
)

type Auth struct {
}

// Business logic core
func (auth *Auth) LoginCore(ctx context.Context, username, password string) (*core.OAuthResponse, error) {
	if strings.Contains(username, "/") {
		parts := strings.SplitN(username, "/", 2)
		tenant, user := parts[0], parts[1]
		oauth, err := core.Services.AuthSvc.Login(tenant, ctx, user, password)
		if err != nil || oauth == nil {
			return nil, wx.NewUnauthorizedError()
		}
		return oauth, nil
	}

	oauth, err := core.Services.TenantSvc.Login(ctx, username, password)
	if err != nil || oauth == nil {
		return nil, wx.NewUnauthorizedError()
	}
	return oauth, nil
}
func (auth *Auth) Login(handler wx.Handler, data *wx.Form[struct {
	Username string
	Password string
}]) (*core.OAuthResponse, error) {
	return auth.LoginCore(handler().Req.Context(), data.Data.Username, data.Data.Password)

}

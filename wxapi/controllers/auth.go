package controllers

import (
	"core"
	"strings"

	"github.com/vn-go/wx"
)

type Auth struct {
}

func (auth *Auth) Login(handler wx.Handler, data *wx.Form[struct {
	Username string
	Password string
}]) (*core.OAuthResponse, error) {
	if strings.Contains(data.Data.Username, "/") {
		tenant := strings.Split(data.Data.Username, "/")[0]
		username := strings.Split(data.Data.Username, "/")[1]
		oauth, err := core.Services.AuthSvc.Login(tenant, handler().Req.Context(), username, data.Data.Password)
		if err != nil {
			return nil, wx.NewUnauthorizedError()
		} else {
			if oauth == nil {
				return nil, wx.NewUnauthorizedError()
			}
			return oauth, nil
		}
	} else {
		oauth, err := core.Services.TenantSvc.Login(handler().Req.Context(), data.Data.Username, data.Data.Password)
		if err != nil {
			return nil, err
		} else {
			if oauth == nil {
				return nil, wx.NewUnauthorizedError()
			}
			return oauth, nil
		}
	}

}

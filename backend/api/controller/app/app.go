package app

import (
	"apicore/controller/base"

	"github.com/vn-go/wx"
)

type App struct {
	base.AuthBase
}

func (a *App) ApiDiscovery(h wx.Handler, apiList []string) (any, error) {
	user := a.Authenticate.Data
	if !user.IsSysAdmin {
		return nil, nil
	}
	return a.Svc.AppSvc.ApiDiscovery(apiList)
}

package controllers

import (
	"core"

	"github.com/vn-go/wx"
)

type ViewManager struct {
	BaseAuthController
}

func (vm *ViewManager) GetListOfViews(h wx.Handler) (data any, err error) {
	return core.Services.ViewSvc.GetListOfViews(h().Req.Context(), vm.Data)
}

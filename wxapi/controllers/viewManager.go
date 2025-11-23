package controllers

import (
	"core"

	"github.com/vn-go/wx"
)

type ViewManager struct {
	BaseAuthController
}

func (vm *ViewManager) GetListOfViews(h wx.Handler) (data any, err error) {

	data, err = core.Services.ViewSvc.GetListOfViews(h().Req.Context(), vm.Data)
	if err != nil {
		return nil, &wx.BadRequestError{
			Message: err.Error(),
		}
	}
	return data, nil
}

func (vm *ViewManager) ApiDiscovery(h wx.Handler, data struct {
	ApiList []string `json:"apiList"`
}) (any, error) {
	return core.Services.ViewSvc.ApiDiscovery(h().Req.Context(), vm.Authenticate.Data, data.ApiList)

}

/*
get list of api by view path
*/
func (vm *ViewManager) GetListApiOfView(h wx.Handler, data struct {
	ViewPath string `json:"viewPath"`
}) (any, error) {
	return core.Services.ViewSvc.GetListApiOfView(h().Req.Context(), vm.Authenticate.Data, data.ViewPath)

}
func (vm *ViewManager) ApiDelete(h wx.Handler, data struct {
	ApiPath string `json:"apiPath"`
}) (any, error) {
	return core.Services.ViewSvc.DeleteApi(h().Req.Context(), vm.Authenticate.Data, data.ApiPath)

}

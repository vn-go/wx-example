package controllers

import (
	"core"

	"github.com/vn-go/dx"
	"github.com/vn-go/wx"
)

type DataSource struct {
	BaseAuthController
}

func (ds *DataSource) Get(h wx.Handler, data struct {
	Name     string `json:"name" `
	Selector string `json:"fields"`
	Filter   string `json:"filter" check:"range:[0:300]"`
}) (any, error) {
	//dx.Options.ShowSql = true
	ret, err := core.Services.DataSvc.GetList(h().Req.Context(), ds.Authenticate.Data, data.Name, data.Filter, data.Selector)
	if exprErr := dx.Errors.IsExpressionError(err); exprErr != nil {
		return nil, wx.Errors.NewHttpError(wx.ErrUnprocessableEntity, exprErr.Error())
	}
	return ret, err

}

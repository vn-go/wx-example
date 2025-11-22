package controllers

import (
	"core"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/vn-go/dx"
	"github.com/vn-go/wx"
)

// DataSource 控制器 controller
type DataSource struct {
	BaseAuthController // require authentication for all actions
}

// Get datasource
func (ds *DataSource) Get(h wx.Handler, data struct {
	// Datasource name, (a certain datasource, has registed at backend)
	Name string `json:"name" `
	// Selector, fields to select, separated by comma, list all fields not allow empty or "*"
	// support function: sum, avg, max, min, count, if(condition-1, trueValue-1, ...,condition-n, falseValue-n), concat(field1, field2,...)
	// example: "id,name,age,sum(age),if(age>18, 'adult',age>13 and age<18, 'teenager', 'unkonwn')"
	Selector string `json:"fields"`
	// filter condition
	// Exampe: "id=1 and name='abc' and sum(salary)>10000"
	Filter string `json:"filter" check:"range:[0:300]"`
}) (any, error) {
	//dx.Options.ShowSql = true

	ret, err := core.Services.DataSvc.GetList(h().Req.Context(), ds.Authenticate.Data, data.Name, data.Filter, data.Selector)
	if exprErr := dx.Errors.IsExpressionError(err); exprErr != nil {
		return nil, wx.Errors.NewHttpError(wx.ErrUnprocessableEntity, exprErr.Error())
	}
	return ret, err

}

func (ds *DataSource) GetQuery(h wx.Handler, data struct {
	Name     string `json:"name" `
	Selector string `json:"fields"`
	Filter   string `json:"filter" check:"range:[0:300]"`
}) (any, error) {
	//dx.Options.ShowSql = true

	ret, err := core.Services.DataSvc.GetSQL(h().Req.Context(), ds.Authenticate.Data, data.Name, data.Filter, data.Selector)
	if exprErr := dx.Errors.IsExpressionError(err); exprErr != nil {
		return nil, wx.Errors.NewHttpError(wx.ErrUnprocessableEntity, exprErr.Error())
	}
	return ret, err

}

var cacheDsl sync.Map

func (ds *DataSource) RegisterDataSource(h wx.Handler, data struct {
	Dsl string `json:"dsl"`
}) (any, error) {
	if ret, ok := cacheDsl.Load(strings.ToLower(data.Dsl)); ok {
		return struct {
			Id string `json:"id"`
		}{
			Id: ret.(string),
		}, nil
	} else {
		retId := uuid.NewString()
		cacheDsl.Store(strings.ToLower(data.Dsl), retId)
		return struct {
			Id string `json:"id"`
		}{
			Id: ret.(string),
		}, nil
	}
}

package controllers

import (
	"core"

	"github.com/vn-go/wx"
)

type Dataset struct {
	BaseAuthController // require authentication for all actions
}
type DatasetInfo struct {
	Source    string `json:"source"`
	Fields    string `json:"fields"`
	Filter    string `json:"filter"`
	PageSize  uint64 `json:"limit"`
	PageIndex uint64 `json:"offset"`
}

func (c *Dataset) List(h wx.Handler, data DatasetInfo) (any, error) {
	ret, err := core.Services.Dataset.Execute(h().Req.Context(), c.Data, data.Source, data.Fields, data.Filter, "", data.PageSize, data.PageIndex)
	if err != nil {
		return nil, wx.Errors.NewHttpError(wx.ErrBadRequest, err.Error())
	}
	return ret, nil
}

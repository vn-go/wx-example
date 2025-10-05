package controllers

import "github.com/vn-go/wx"

type Pure struct {
	//BaseAuthController
}

func (ds *Pure) Get(h wx.Handler) (any, error) {
	return "OK", nil
}

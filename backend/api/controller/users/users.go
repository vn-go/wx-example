package users

import (
	"apicore/controller/base"

	"github.com/vn-go/wx"
)

type Users struct {
	base.AuthBase
}

func (u *Users) List(h wx.Handler) (any, error) {
	return u.Svc.AccSvc.ListUsers(h().Req.Context(), u.Authenticate.Data)
}
func (u *Users) Me(h wx.Handler) (any, error) {
	req := h().Req

	return u.Svc.AccSvc.CurrentUserProfile(req.Context(), u.Authenticate.Data)
}

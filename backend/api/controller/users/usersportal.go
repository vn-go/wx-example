package users

import (
	"core"

	"github.com/vn-go/wx"
)

type UserPortal struct {
}

func (u *UserPortal) GetList(h wx.Handler) (any, error) {
	return core.Services.Db.DslToArray("sysUsers()")
}

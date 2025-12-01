package users

import "github.com/vn-go/wx"

func (u *Users) Roles(h wx.Handler) (any, error) {
	db, err := u.Svc.TenantSvc.GetDb(u.Data.Tenant)
	if err != nil {
		return nil, err
	}
	return db.DslToArray("sysRoles(id,code,name),sort(name asc)")
}

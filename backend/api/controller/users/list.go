package users

import "github.com/vn-go/wx"

type Pager struct {
	Index int `json:"index"`
	Size  int `json:"size" check:"[1:50]"`
}

func (u *Users) List(h wx.Handler) (any, error) {
	db, err := u.Svc.TenantSvc.GetDb(u.Data.Tenant)
	if err != nil {
		return nil, err
	}
	ret, err := db.DslToArray(`
			sysUsers(), 
			sysRoles(code RoleCode, name RoleName),
			from(sysUsers  u, sysRoles  r,left(u.roleId=r.Id)),
			sort(u.createdOn desc)
	`)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

/*
//C:\Users\MSI CYBORG\AppData\Local\Temp
$SourcePath = "C:\Users\MSI CYBORG\AppData\Local\JetBrains\Fleet\cache"
$TargetPath = "D:\JetBrainsCache\Fleet\cache"
cmd /c mklink /J $SourcePath $TargetPath
*/

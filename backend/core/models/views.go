package models

import "github.com/vn-go/dx"

type SysViews struct {
	ViewPath string `db:"size(250);uk"`
	Title    string `db:"size(250);uk"`
	baseModel
}
type SysViewRoles struct {
	ViewId string `db:"pk;size(36);" json:"viewId"`
	RoleId string `db:"pk;size(36);" json:"roleId"`
}

func init() {
	dx.AddForeignKey[SysViewRoles]("viewId", &SysViews{}, "Id", nil)
	dx.AddForeignKey[SysViewRoles]("RoleId", &SysRoles{}, "Id", nil)
}

package models

import "github.com/vn-go/dx"

type SysApi struct {
	Enpoint string `db:"size(250);uk"`
	Title   string `db:"size(250);uk"`
	baseModel
}
type SysViewApi struct {
	ApiId  string `db:"pk;size(36); default:uuid()" json:"apiId"`
	ViewId string `db:"pk;size(36); default:uuid()" json:"viewId"`
}

func init() {
	dx.AddForeignKey[SysViewApi]("ApiId", &SysApi{}, "Id", nil)
	dx.AddForeignKey[SysViewApi]("ViewId", &SysViews{}, "Id", nil)
}

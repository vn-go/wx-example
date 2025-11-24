package models

import "github.com/vn-go/dx"

type SysUsers struct {
	baseModel
	Username     string  `db:"size(50);uk" json:"username"`
	DisplayName  *string `db:"size(50); default:''" json:"displayName"`
	Email        string  `db:"size(50)" json:"email"`
	HashPassword string  `db:"size(150);" json:"-"`
	RoleId       *string `db:"size(36);idx" json:"roleId"`
	IsSysAdmin   bool    `json:"isSysAdmin"`
}

func init() {
	dx.AddForeignKey[SysUsers]("roleId", &SysRoles{}, "id", &dx.FkOpt{
		OnDelete: true,
		OnUpdate: true,
	})
}

package models

import "github.com/vn-go/dx"

type SysRefreshToken struct {
	Username string `db:"pk;size(50)"`
	// Token    string `db:"size(450)"`
	RefreshToken string `db:"size(36)"`
	baseInfo
}

func init() {
	dx.AddModels(&SysRefreshToken{})
}

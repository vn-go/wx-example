package models

import (
	"github.com/vn-go/dx"
)

type SysApp struct {
	Code string `db:"pk;size(50)" json:"code"`
	Name string `db:"size(250)" json:"name"`

	baseInfo
	SecretKey string `db:"size(250)" json:"-"`
}

func init() {
	dx.AddModels(&SysApp{})
}

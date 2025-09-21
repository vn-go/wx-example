package models

import (
	"time"

	"github.com/vn-go/dx"
)

type Tenant struct {
	Id       uint64  `db:"pk;auto"`
	TenantId string  `db:"size:36;default:uuid()"`
	Name     string  `db:"size:50;uk"`
	DbName   string  `db:"size:50;uk"`
	Email    *string `db:"size:50;uk"`
	// ix -> index
	CreatedOn  time.Time `db:"ix;default:now()"`
	ModifiedOn *time.Time
	//df means default value
	IsActive        bool       `db:"default:true"`
	LatestLoginFail *time.Time `db:"ix"`
	LatestLogin     *time.Time `db:"ix"`
	ShareSecret     string     `db:"size:500;default:''"`
}
type RefreshToken struct {
	Token  string `db:"size:36"`
	UserId uint64 `db:"pk;size:36"`
}

func init() {
	dx.AddModels(&Tenant{}, &RefreshToken{})
}

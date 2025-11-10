package models

import (
	"time"

	"github.com/vn-go/dx"
)

type Tenant struct {
	Id       uint64  `db:"pk;auto" json:"id"`
	TenantId string  `db:"size:36;default:uuid()"  json:"tenantId"`
	Name     string  `db:"size:50;uk"  json:"name"`
	DbName   string  `db:"size:50;uk"  json:"-"`
	Email    *string `db:"size:50;uk"  json:"email"`
	// ix -> index
	CreatedOn  time.Time  `db:"ix;default:now()"  json:"createdOn"`
	ModifiedOn *time.Time `json:"modifiedOn"`
	//df means default value
	IsActive        bool       `db:"default:true"  json:"isActive"`
	LatestLoginFail *time.Time `db:"ix"  json:"latestLoginFail"`
	LatestLogin     *time.Time `db:"ix"  json:"latestLogin"`
	ShareSecret     string     `db:"size:500;default:''"  json:"-"`
}
type RefreshToken struct {
	Token  string `db:"size:36" json:"-"`
	UserId string `db:"pk;size:36" json:"-"`
}
type App struct {
	Name              string     `db:"pk;size:255"  json:"name"`
	ShareSecret       string     `db:"size:500"  json:"-"`
	AdminUsername     string     `db:"size:100;default:''" json:"-"`
	AdminHashPassword string     `db:"pk;size:200;default:''" json:"-"`
	CreatedOn         time.Time  `db:"ix;default:now()" json:"createdOn"`
	ModifiedOn        *time.Time `json:"modifiedOn"`
	MasterSecret      string     `db:"size:500;default:''" json:"-"`
}

func (r *RefreshToken) Table() string {
	return "sys_refresh_tokens"
}
func (r *Tenant) Table() string {
	return "sys_tenants"
}
func (r *App) Table() string {
	return "sys_apps"
}
func init() {
	dx.AddModels(&Tenant{}, &RefreshToken{}, &App{})
}

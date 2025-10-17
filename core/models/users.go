package models

import (
	"time"

	"github.com/vn-go/dx"
)

type User struct {
	Id uint64 `db:"pk;auto" json:"id"`
	// UserId is a text of guid
	// uk: uniquekey
	// size:36 -> length of field
	UserId       string `db:"uk;size:36;default:uuid()" json:"userId"`
	Username     string `db:"size:50;uk" json:"username"`
	HashPassword string `db:"size:200" json:"-"`
	// any field is ptr -> allow null
	Email *string `db:"size:50;uk" json:"email"`
	// ix -> index
	CreatedOn  time.Time `db:"ix;default:now()" json:"createdOn"`
	ModifiedOn *time.Time
	//df means default value
	IsActive               bool       `db:"default:true" json:"isActive"`
	LatestLoginFail        *time.Time `db:"ix" json:"latestLoginFail"`
	LatestLogin            *time.Time `db:"ix"  json:"latestLogin"`
	RoleCode               *string    `db:"size:50;ix" json:"roleCode"`
	LastTimeChangePassword *time.Time `json:"lastTimeChangePassword"`
	IsTenantAdmin          *bool      `json:"-"`
	RoleId                 *uint64    `json:"roleId"`
	IsSysAdmin             bool       `db:"default:false" json:"isSysAdmin"`
	CreatedBy              string     `db:"size:50;default:'admin'"  json:"createdBy"`
}
type Role struct {
	Id uint64 `db:"pk;auto" json:"id"`
	// UserId is a text of guid
	// uk: uniquekey
	// size:36 -> length of field
	RoleId      string     `db:"uk;size:36;default:uuid()" json:"roleId"`
	Code        string     `db:"size:50;uk" json:"code"`
	Name        string     `db:"size:50;uk" json:"name"`
	Description *string    `db:"size:200" json:"description"`
	CreatedOn   time.Time  `db:"ix;default:now()" json:"createdOn"`
	ModifiedOn  *time.Time `json:"modifiedOm"`
	CreatedBy   string     `db:"size:50" json:"createdBy"`
	IsActive    bool       `db:"default:true" json:"isActive"`
}

func (u *Role) Table() string {
	return "sys_roles"
}
func (u *User) Table() string {
	return "sys_users"
}

type TrackFilter struct {
	Id uint64 `db:"pk;auto" json:"id"`
	// UserId is a text of guid
	// uk: uniquekey
	// size:36 -> length of field
	DsName   string `db:"size:50"`
	Filter   string `db:"size:2000"`
	Selector string `db:"size:2000"`
}

func init() {
	//rehgister User
	dx.AddModels(&User{}, &Role{}, &TrackFilter{})
	dx.AddForeignKey[User]("RoleId", &Role{}, "Id", &dx.FkOpt{
		OnDelete: false,
		OnUpdate: false,
	})
}

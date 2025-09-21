package models

import (
	"time"

	"github.com/vn-go/dx"
)

type User struct {
	Id uint64 `db:"pk;auto"`
	// UserId is a text of guid
	// uk: uniquekey
	// size:36 -> length of field
	UserId       string `db:"uk;size:36;default:uuid()"`
	Username     string `db:"size:50;uk"`
	HashPassword string `db:"size:200"`
	// any field is ptr -> allow null
	Email *string `db:"size:50;uk"`
	// ix -> index
	CreatedOn  time.Time `db:"ix;default:now()"`
	ModifiedOn *time.Time
	//df means default value
	IsActive        bool       `db:"default:true"`
	LatestLoginFail *time.Time `db:"ix"`
	LatestLogin     *time.Time `db:"ix"`
	RoleCode        *string    `db:"size:50;ix"`
}

func init() {
	//rehgister User
	dx.AddModels(&User{})
}

package models

import "github.com/vn-go/xdb"

type User struct {
	xdb.Model[User]
	ID           uint64 `db:"auto"`
	UserId       string `db:"df:uuid();size:36"`
	Username     string `db:"uk;size:50"`
	HashPassword string `db:"size:200"`
}

func init() {
	xdb.ModelRegistry.Add(&User{})
}

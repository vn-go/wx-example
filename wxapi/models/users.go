package models

import (
	"time"

	
)

type User struct {
	xdb.Model[User]
	ID           uint64    `db:"auto;pk"`
	UserId       string    `db:"df:uuid();size:36;uk"`
	Username     string    `db:"uk;size:50"`
	HashPassword string    `db:"size:200"`
	Email        string    `db:"size:50;uk"`
	CreatedOn    time.Time `db:"df:now()"`
	ModifiedOn   *time.Time
}

func init() {
	xdb.ModelRegistry.Add(&User{})
}

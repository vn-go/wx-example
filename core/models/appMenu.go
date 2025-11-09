package models

import (
	"time"

	"github.com/vn-go/dx"
)

type AppMenu struct {
	// Id of menu
	ID        uint64  `db:"pk;auto"`
	ParentId  *uint64 `db:"index"`
	Title     string  `db:"size:255;unique"`
	IdPaths   string  `db:"size:500"`
	Icon      string  `db:"size:50"`
	ViewPath  string  `db:"size:255;unique"`
	CreatedBy string  `db:"size:50"`
	CreatedOn time.Time
	UpdatedBy string `db:"size:50"`
	UpdatedOn *time.Time
}

func (m *AppMenu) Table() string {
	return "sys_app_menu"
}
func init() {
	dx.AddForeignKey[AppMenu]("ParentId", &AppMenu{}, "ID", &dx.FkOpt{
		OnDelete: true,
	})
}

package app

import (
	"core/models"

	"github.com/vn-go/dx"
)

func (app *AppService) createDefaultAccount() error {
	cfg := app.cfgSvc.Get()
	accItem, err := dx.NewDTO[models.SysUsers]()
	if err != nil {
		return err
	}
	accItem.Username = cfg.DefaultAuth.Username
	accItem.HashPassword, err = app.pwdSvc.HashPassword(accItem.Username, cfg.DefaultAuth.Password)
	if err != nil {
		return err

	}
	accItem.IsSysAdmin = true
	err = app.db.Insert(accItem)
	dbErr := dx.Errors.IsDbError(err)

	if dbErr == nil || dbErr.IsDuplicateEntryError() {
		return nil
	}
	return dbErr

}

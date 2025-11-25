package app

import (
	"core/models"
	"time"

	"github.com/vn-go/dx"
)

func (app *AppService) InitSingleTenant() error {
	appItem, err := dx.NewDTO[models.SysApp]()
	if err != nil {
		return err
	}
	appItem.Code = "App"
	appItem.Name = "Application"
	appItem.SecretKey, err = app.jwtSvc.GenerateSecret(0)
	if err != nil {
		return err
	}
	appItem.CreatedBy = "application"
	appItem.CreatedOn = time.Now().UTC()
	appItem.AesKey, err = app.GenerateRandomKey()
	if err != nil {
		return err
	}
	err = app.db.Insert(appItem)
	dbErr := dx.Errors.IsDbError(err)
	if dbErr != nil {
		if !dbErr.IsDuplicateEntryError() {
			return dbErr
		} else {
			if appItem.Name == "" {
				appItem.Name = "Application"
			}
			if appItem.SecretKey == "" {
				appItem.SecretKey, err = app.jwtSvc.GenerateSecret(0)
				if err != nil {
					return err
				}

			}
			if appItem.AesKey == "" {
				appItem.AesKey, err = app.GenerateRandomKey()
				if err != nil {
					return err
				}
			}
			rs := app.db.Update(appItem)
			if rs.Error != nil {
				if dbErr = dx.Errors.IsDbError(rs.Error); dbErr != nil {
					if !dbErr.IsDuplicateEntryError() {
						return dbErr
					}
				} else {
					return rs.Error
				}
			}
		}

	}
	return nil
}

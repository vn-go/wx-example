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
	appItem.CreatedBy = "application"
	appItem.CreatedOn = time.Now().UTC()
	if err != nil {
		return err
	}
	err = app.db.Insert(appItem)
	dbErr := dx.Errors.IsDbError(err)
	if dbErr != nil {
		if !dbErr.IsDuplicateEntryError() {
			return dbErr
		} else {
			// err = app.db.First(appItem, "code=?", appItem.Code)
			// if err != nil {
			// 	return err
			// }
			if appItem.SecretKey == "" {
				appItem.SecretKey, err = app.jwtSvc.GenerateSecret(0)
				if err != nil {
					return err
				}
				if appItem.Name == "" {
					appItem.Name = "Application"
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

	}
	return nil
}

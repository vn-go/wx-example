package acc

import (
	"context"
	"core/models"
	"time"

	"github.com/google/uuid"
	"github.com/vn-go/dx"
)

func (acc *AccService) CreateRefreshToken(ctx context.Context, tenant, username string, tk string) (string, error) {
	db, err := acc.tenantSvc.GetDb(tenant)
	if err != nil {
		return "", err
	}
	refreshItem, err := dx.NewDTO[models.SysRefreshToken]()
	if err != nil {
		return "", err
	}
	refreshItem.Username = username
	refreshItem.CreatedBy = "application"
	refreshItem.CreatedOn = time.Now().UTC()
	refreshItem.RefreshToken = uuid.NewString()
	err = db.InsertWithContext(ctx, refreshItem)
	dbErr := dx.Errors.IsDbError(err)
	if dbErr != nil {
		if dbErr.IsDuplicateEntryError() {
			rs := db.UpdateWithContext(ctx, refreshItem)
			if rs.Error != nil {
				return "", rs.Error
			}
		}
	}
	return refreshItem.RefreshToken, nil
}

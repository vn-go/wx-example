package acc

import (
	"context"
	"core/models"
	"time"

	"github.com/vn-go/dx"
)

func (acc *AccService) NewUser(ctx context.Context, tenant string, user *UserInfo) error {
	db, err := acc.tenantSvc.GetDb(tenant)
	if err != nil {
		return err
	}
	userItem, err := dx.NewDTO[models.SysUsers]()
	if err != nil {
		return err
	}
	userItem.CreatedBy = "admin"
	userItem.CreatedOn = time.Now().UTC()
	userItem.Email = user.Email
	userItem.Username = user.Username
	userItem.HashPassword, err = acc.pwdSvc.HashPassword(user.Username, user.Password)
	userItem.DisplayName = &user.DisplayName
	if err != nil {
		return err
	}
	err = db.Insert(userItem)
	return err
}

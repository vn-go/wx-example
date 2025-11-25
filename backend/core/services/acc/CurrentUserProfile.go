package acc

import (
	//"apicore/services/jwt"
	"context"
	"core/models"
	"core/services/jwt"
)

func (acc *AccService) CurrentUserProfile(context context.Context, user *jwt.Indentifier) (any, error) {
	db, err := acc.tenantSvc.GetDb(user.Tenant)
	if err != nil {
		return nil, err
	}
	var userProfile = &models.SysUsers{}
	err = db.First(userProfile, "id = ?", user.UserId)
	return userProfile, err
	// TODO: implement

}

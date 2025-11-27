package users

import (
	"apicore/controller/base"
	"core"
	"core/models"
	"time"

	"github.com/vn-go/wx"
)

type Users struct {
	base.AuthBase
}

func (u *Users) Me(h wx.Handler) (any, error) {
	req := h().Req

	return u.Svc.AccSvc.CurrentUserProfile(req.Context(), u.Authenticate.Data)
}

func (u *Users) GetItem(h wx.Handler, data struct {
	UserId string `json:"userId"`
}) (*editUser, error) {
	db, err := u.Svc.TenantSvc.GetDb(u.Authenticate.Data.Tenant)
	if err != nil {
		return nil, err
	}
	ret := &editUser{
		Data: models.SysUsers{},
	}
	err = db.First(&ret.Data, "id = ?", data.UserId)
	if err != nil {
		return nil, u.ParseError(err)
	}
	err = u.Svc.DataSvc.SignData(h().Req.Context(), u.Authenticate.Data, ret)
	if err != nil {
		return nil, u.ParseError(err)
	}
	return ret, nil
}

type editUser core.DataContract[models.SysUsers, struct {
	Id           string     // can not modify by client
	Username     string     // can not modify by client
	CreatedBy    string     // can not modify by client
	CreatedOn    time.Time  // can not modify by client
	ModifiedBy   *string    // can not modify by client
	ModifiedOn   *time.Time // can not modify by client
	HashPassword string     // can not modify by client
}]

func (u *Users) Update(h wx.Handler, data *editUser) (*editUser, error) {
	//verify data from client and also extract key field in token then update
	// to u.Data
	// exmaple:
	// if client modify email, then update email in SysUsers table error will be returned
	err := u.Svc.DataSvc.Verify(h().Req.Context(), u.Data, data)
	if err != nil {
		return nil, u.ParseError(err)
	}
	// get tenant db of current user login
	db, err := u.Svc.TenantSvc.GetDb(u.Authenticate.Data.Tenant)
	if err != nil {
		return nil, err
	}
	// update SysUsers table afer all verification
	rs := db.UpdateWithContext(h().Req.Context(), data.Data)
	if rs.Error != nil {
		return nil, u.ParseError(rs.Error)
	}
	return data, nil
}

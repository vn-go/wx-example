package roles

import (
	"apicore/controller/base"
	"core"
	"core/models"

	"github.com/vn-go/dx"
	"github.com/vn-go/wx"
)

type Roles struct {
	base.AuthBase
}

func (r *Roles) List(h wx.Handler) (data any, err error) {
	db, err := r.AuthBase.GetUserDb()
	if err != nil {
		return
	}
	data, err = db.DslToArray(`
					SysRoles( 		 //select id,name,... in sys_roles			
							id,name,
							code,
							description,
							createdOn,
							createdBy,
							modifiedOn,
							modifiedBy),
					SysUsers( count(id) TotalUsers),
							from(SysRoles r, SysUsers u,left(r.Id=u.RoleId)), // left join SysUsers on SysRoles.Id=SysUsers.RoleId
							sort(createdOn desc)`)
	return
}

type RoleNewItemLock struct {
	Id        string
	CreatedBy string
}
type DataRow[TModel any] struct {
	Data   TModel
	Status int
}

func (r *Roles) NewItem(h wx.Handler) (data *core.DataContract[models.SysRoles, RoleNewItemLock], err error) {
	ctx := h()
	roleItem, err := dx.NewDTO[models.SysRoles]()
	if err != nil {
		return nil, err
	}
	roleItem.CreatedBy = r.Data.Username
	data = &core.DataContract[models.SysRoles, RoleNewItemLock]{
		Data:   *roleItem,
		Status: "new",
	}
	err = r.Svc.DataSvc.SignData(ctx.Req.Context(), r.Data, data)
	if err != nil {
		return nil, err
	}
	return
}
func (r *Roles) Update(h wx.Handler, data *core.DataContract[models.SysRoles, RoleNewItemLock]) (any, error) {
	ctx := h()
	db, err := r.AuthBase.GetUserDb()
	if err != nil {
		return nil, err
	}
	err = r.Svc.DataSvc.Verify(ctx.Req.Context(), r.Data, data)
	if err != nil {
		return nil, r.ParseError(err)
	}
	if data.Status == "new" {
		err = db.InsertWithContext(ctx.Req.Context(), &data.Data)
		if err != nil {
			return nil, r.ParseError(err)
		}
		return data, nil
	} else {
		rs := db.UpdateWithContext(ctx.Req.Context(), &data.Data)
		if rs.Error != nil {
			return nil, r.ParseError(rs.Error)
		}
	}
	return data, nil
}
func (r *Roles) Item(h wx.Handler, id string) (data *core.DataContract[models.SysRoles, RoleNewItemLock], err error) {
	ctx := h()
	db, err := r.AuthBase.GetUserDb()
	if err != nil {
		return
	}
	roleItem := models.SysRoles{}
	err = db.First(&roleItem, "id=?", id)
	if err != nil {
		return nil, r.ParseError(err)
	}
	data = &core.DataContract[models.SysRoles, RoleNewItemLock]{
		Data:   roleItem,
		Status: "edit",
	}
	err = r.Svc.DataSvc.SignData(ctx.Req.Context(), r.Data, data)
	if err != nil {
		return nil, err
	}
	return
}

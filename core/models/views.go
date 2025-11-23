package models

import (
	"time"

	"github.com/vn-go/dx"
)

type UIView struct {
	Id uint64 `db:"pk;auto"`
	//sha256 of view viewPath
	ViewId      string `db:"size:64;unique"`
	ViewPath    string `db:"size:255;unique"`
	Title       string `db:"size:250;default:''"`
	Description string `db:"size:255"`
	CreatedOn   time.Time
	ModifiedOn  *time.Time
	CreatedBy   string  `db:"size:50"`
	ModifiedBy  *string `db:"size:50"`
}

func (*UIView) Table() string {
	return "sys_view"
}

type RoleUiView struct {
	Id          uint64 `db:"pk;auto"`
	RoleId      uint64
	ViewId      uint64
	AllowInsert bool
	AllowUpdate bool
	AllowDelete bool
	CreatedBy   string `db:"size:50"`
}

func (*RoleUiView) Table() string {
	return "sys_role_view"
}

type UIDataSource struct {
	Id uint64 `db:"pk;auto"`
	//sha256 of Dsl
	DataSourceId string `db:"size:64;unique"`
	Dsl          string
	CreatedOn    time.Time
	ModifiedOn   *time.Time
	CreatedBy    string  `db:"size:50"`
	ModifiedBy   *string `db:"size:50"`
}

func (*UIDataSource) Table() string {
	return "sys_datasource"
}

type UIViewApi struct {
	Id          uint64 `db:"pk;auto" json:"-"`
	ViewId      uint64 `db:"unique:view_api_uix"`
	ApiId       uint64 `db:"unique:view_api_uix"`
	Title       string `db:"size:250;default:''"`
	Description string
	CreatedOn   time.Time
	ModifiedOn  *time.Time
	CreatedBy   string `db:"size:50"`
}

func (*UIViewApi) Table() string {
	return "sys_view_api"
}

type Api struct {
	Id          uint64 `db:"pk;auto"`
	ApiPath     string `db:"size:255;unique"`
	Title       string `db:"size:250;default:''"`
	Description string
	CreatedOn   time.Time
	CreatedBy   string `db:"size:50;idx;default:''"`
	ModifiedOn  *time.Time
}

func (*Api) Table() string {
	return "sys_api"
}

type UIDataSourceDetail struct {
	Id           uint64 `db:"pk;auto"`
	DataSourceId uint64
	Field        string
	Expr         string
}

func (*UIDataSourceDetail) Table() string {
	return "sys_datasource_detail"
}

type UIViewDataSource struct {
	Id           uint64 `db:"pk;auto"`
	ViewId       uint64 `db:"unique:view_dsid_uix"`
	DataSourceId uint64 `db:"unique:view_dsid_uix"`
}

func (*UIViewDataSource) Table() string {
	return "sys_view_datasource"
}
func init() {
	dx.AddModels(&UIView{}, &UIDataSource{}, &UIViewDataSource{})
	dx.AddForeignKey[UIViewDataSource]("ViewId", &UIView{}, "Id", &dx.FkOpt{
		OnDelete: true,
	})
	dx.AddForeignKey[UIViewDataSource]("ViewId", &UIView{}, "Id", &dx.FkOpt{
		OnDelete: true,
	})
	dx.AddForeignKey[UIDataSourceDetail]("DataSourceId", &UIDataSource{}, "Id", &dx.FkOpt{
		OnDelete: true,
	})
	dx.AddForeignKey[UIViewDataSource]("DataSourceId", &UIDataSource{}, "Id", &dx.FkOpt{
		OnDelete: true,
	})
	dx.AddForeignKey[RoleUiView]("RoleId", &Role{}, "Id", &dx.FkOpt{
		OnDelete: true,
	})
	dx.AddForeignKey[RoleUiView]("ViewId", &UIView{}, "Id", &dx.FkOpt{
		OnDelete: true,
	})
	dx.AddForeignKey[UIViewApi]("ApiId", &Api{}, "Id", &dx.FkOpt{
		OnDelete: true,
	})
	dx.AddForeignKey[UIViewApi]("ViewId", &UIView{}, "Id", &dx.FkOpt{
		OnDelete: true,
	})
}

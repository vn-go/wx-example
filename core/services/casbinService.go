package services

import (
	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vn-go/dx"
)

func NewCasbinEnforcer(db *dx.DB) (*casbin.Enforcer, error) {

	a, err := sqladapter.NewAdapter(db.DB.DB, db.DriverName, "casbin_rule")
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer("../models/model.conf", a)
	if err != nil {
		return nil, err
	}
	return e, nil
}

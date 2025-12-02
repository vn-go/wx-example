package controller

import (
	"apicore/controller/app"
	"apicore/controller/auth"
	"apicore/controller/media"
	"apicore/controller/roles"
	"apicore/controller/users"

	"github.com/vn-go/wx"
)

func InitRoutes() {
	err := wx.Routes("system",
		&auth.Auth{},
		&users.Users{},
		&app.App{},
		&roles.Roles{},
		&users.UserPortal{},
		&media.Media{},
	)
	if err != nil {
		panic(err)
	}
}

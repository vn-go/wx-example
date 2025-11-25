package controller

import (
	"apicore/controller/auth"
	"apicore/controller/users"

	"github.com/vn-go/wx"
)

func InitRoutes() {
	err := wx.Routes("/api",
		&auth.Auth{},
		&users.Users{},
	)
	if err != nil {
		panic(err)
	}
}

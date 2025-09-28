package routes

import (
	"wxapi/controllers"

	"github.com/vn-go/wx"
)

func InitRoute() {
	err := wx.Routes("/api",
		&controllers.Auth{},
		&controllers.Tenant{},
		&controllers.Accounts{})
	if err != nil {
		panic(err)
	}
	// err := wx.Routes("/api", &controllers.Auth{}, &controllers.Tenant{}, &controllers.Accounts{})
	// if err != nil {
	// 	panic(err)
	// }
}

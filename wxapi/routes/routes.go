package routes

import (
	"wxapi/controllers"

	"github.com/vn-go/wx"
)

func InitRoute() {
	controllers.GetListOfRoles()
	// acc all controllers to routes
	err := wx.Routes("/api",

		&controllers.Auth{},
		&controllers.Tenant{},
		&controllers.Accounts{},
		&controllers.DataSource{},
		&controllers.Pure{},
	)
	if err != nil {
		panic(err)
	}
	// err := wx.Routes("/api", &controllers.Auth{}, &controllers.Tenant{}, &controllers.Accounts{})
	// if err != nil {
	// 	panic(err)
	// }
}

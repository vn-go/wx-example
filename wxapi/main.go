package main

import (
	"wxapi/controllers"
	"wxapi/routes"

	"github.com/vn-go/wx"
)

type Hello struct {
}

func (hello *Hello) Hx(ctx *struct {
	wx.Handler `route:"method:get"`
}) {
	ctx.Handler().Res.Write([]byte("ok"))
}

type DataPost struct {
	Name string
}

func (t *Hello) Create(h wx.Handler, data DataPost) (any, error) {
	//core.Services.TenantSvc.CreateTenant(h().Req.Context(), data.Name, data.Name)
	return struct{}{}, nil
}
func main() {
	routes.InitRoute()
	wx.Routes("/api", &controllers.Auth{}, &controllers.Tenant{})
	wx.Routes("/api", &Hello{})
	server := wx.NewHtttpServer("/api", "8080", "0.0.0.0")
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.OAuth2Password("/api/auth/login")
	swagger.Build()
	server.Start()
}

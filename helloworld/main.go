package main

import (
	"github.com/vn-go/wx"
)

type App struct {
}

func (a *App) Hello(h *struct {
	wx.Handler `route:"method:get"`
}) (any, error) {
	return "Hello, World!", nil
}
func main() {
	wx.Routes("/api", &App{})
	server := wx.NewHtttpServer("/api", "8081", "0.0.0.0")
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.Build()
	server.Start()
}

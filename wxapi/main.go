package main

import "github.com/vn-go/wx"

type Hello struct {
}

func (hello *Hello) Hello(ctx *struct {
	wx.Handler `route:"method:get"`
}) {
	ctx.Handler().Res.Write([]byte("Hello, wx!"))
}
func main() {
	wx.Routes("/api", &Hello{})
	server := wx.NewHtttpServer("/api", "8080", "0.0.0.0")
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.Build()
	server.Start()
}

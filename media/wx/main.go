package main

import (
	"media/controllers/media"
	_ "media/controllers/media"

	"github.com/vn-go/wx"
)

func main() {
	wx.Routes("/api", &media.Media{})
	server := wx.NewHtttpServer("/api", "8080", "0.0.0.0")
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.Build()
	server.Start()
}

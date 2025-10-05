package main

import (
	"fmt"
	"log"
	"media/controllers/media"
	_ "media/controllers/media"
	"net/http"
	_ "net/http/pprof"

	"github.com/vn-go/wx"
)

func main() {
	media.LoadAllFile()
	go func() {
		fmt.Println("pprof running at :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wx.Routes("/api", &media.Media{})
	wx.HandlerGet("file/list-v2", func(controllerInstance *media.Media, ctx *wx.HttpContext[any]) ([]string, error) {
		return controllerInstance.ListAllFolderAndFiles(ctx.GetAbsRootUri())
	})
	server := wx.NewHtttpServer("/api", "8080", "0.0.0.0")
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.Build()
	server.Start()
}

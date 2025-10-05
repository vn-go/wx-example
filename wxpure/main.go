package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vn-go/wx"
)

type Hz struct {
}

func (h *Hz) Check(hd *wx.Handler) (any, error) {
	return "ok", nil
}
func main() {
	go func() {
		fmt.Println("pprof running at :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//wx.Options.IsDebug = core.Services.Config.Debug
	//wx.Options.UsePool = true
	//routes.InitRoute()
	wx.Routes("/api", &Hz{})

	//wx.Routes("/api", &Hello{})
	server := wx.NewHtttpServer("/api", "8080", "0.0.0.0")
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.OAuth2Password("/api/auth/login")
	swagger.Build()
	server.Start()
}

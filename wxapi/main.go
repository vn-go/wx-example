package main

import (
	"core"
	"fmt"
	"log"
	"net/http"
	"wxapi/routes"

	_ "net/http/pprof"

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

func main2() {
	go func() {
		fmt.Println("pprof running at :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	fmt.Println("server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	go func() {
		fmt.Println("pprof running at :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wx.Options.IsDebug = core.Services.Config.Debug
	//wx.Options.UsePool = true
	routes.InitRoute()

	//wx.Routes("/api", &Hello{})
	server := wx.NewHtttpServer("/api", "8080", "0.0.0.0")
	swagger := wx.CreateSwagger(server, "/docs")
	swagger.OAuth2Password("/api/auth/login")
	swagger.Build()
	server.Start()
}

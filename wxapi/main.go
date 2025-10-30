package main

import (
	"core"
	"fmt"
	"log"
	"net/http"
	"time"
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

	// load all controllerd
	routes.InitRoute()

	//new server api
	server := wx.NewHtttpServer("/api", core.Services.Config.Bind.Port, core.Services.Config.Bind.Host)
	server.Middleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		t := time.Now()
		next(w, r) // set middleware if you need
		w.Header().Set("X-Process-Time", fmt.Sprintf("%.2fms", float64(time.Since(t).Nanoseconds())/1e6))
		//fmt.Println("time elapsed:", time.Since(t))
	})
	server.Middleware(wx.MiddlWares.Cors)
	//Create swager if you need
	swagger := wx.CreateSwagger(server, "/docs")
	// Show authenication login in swagger
	swagger.OAuth2Password("/api/auth/login")
	// Build swagger
	swagger.Build()
	// start server
	server.Start()
}

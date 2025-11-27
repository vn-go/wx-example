package main

import (
	"core"
	_ "core"
	"fmt"
	"log"
	"net/http"
	"time"

	"apicore/controller"
	"apicore/middleware"

	"github.com/vn-go/dx"
	"github.com/vn-go/wx"
)

func main() {
	go func() {
		fmt.Println("pprof running at :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	core.Start("./config.yaml")
	defer core.Services.Close()

	wx.Options.IsDebug = core.Services.ConfigSvc.Get().Debug
	wx.OnError(func(err error) {
		fmt.Println(err.Error())
	})
	dx.Options.ShowSql = wx.Options.IsDebug

	// load all controllerd
	controller.InitRoutes()

	//new server api
	server := wx.NewHttpServer(
		"/api",
		core.Services.ConfigSvc.Get().Bind.Port,
		core.Services.ConfigSvc.Get().Bind.Host,
	)

	server.Use(middleware.CorsMiddleware())

	server.Use(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()

		w, r = server.BeforeRequestCompleted("Server-Timing", w, r, func(w http.ResponseWriter, r *http.Request) {
			//elapsed := float64(time.Since(start).Microseconds()) / 1000.0
			durMs := time.Since(start).Seconds() * 1000

			// Format theo chuẩn Server-Timing: <metric-name>;dur=<duration>
			w.Header().Set("Server-Timing", fmt.Sprintf("app;dur=%.2f", durMs))

		})

		next(w, r)
	})
	//Create swager if you need
	swagger := wx.CreateSwagger(server, "/docs")
	// Show authenication login in swagger

	swagger.OAuth2Password("/api/auth/login")
	// Build swagger
	swagger.Build()
	// start server
	server.Start()
}

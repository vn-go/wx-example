package main

import (
	"core"
	"fmt"
	"log"
	"net/http"
	"time"
	"wxapi/routes"

	_ "net/http/pprof"

	"github.com/vn-go/dx"
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
func corsMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var cors = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		// Cho phép tất cả origin (cẩn thận với sản phẩm thật!)
		if len(r.Header["Origin"]) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", r.Header["Origin"][0])
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization,view-path")
		//w.Header().Set("Access-Control-Allow-Origin", "https://frontend.example.com")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Nếu là preflight request (OPTIONS), chỉ phản hồi 200
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Gọi tiếp handler chính
		next.ServeHTTP(w, r)
		//w.Header().Set("X-Process-Time", fmt.Sprintf("%.2fms", float64(time.Since(start).Nanoseconds())/1e6))

	}
	return cors
}
func main() {
	go func() {
		fmt.Println("pprof running at :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wx.Options.IsDebug = core.Services.Config.Debug
	dx.Options.ShowSql = wx.Options.IsDebug

	// load all controllerd
	routes.InitRoute()

	//new server api
	server := wx.NewHttpServer("/api", core.Services.Config.Bind.Port, core.Services.Config.Bind.Host)

	server.Use(corsMiddleware())

	server.Use(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()

		// Đăng ký callback sau khi request hoàn tất
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

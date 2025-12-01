package main

import (
	"apicore/controller"
	"apicore/middleware"
	"core"
	_ "core"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
	_ "net/http/pprof"   // ← thiếu dòng này = 404 100%
	"github.com/vn-go/dx"
	"github.com/vn-go/wx"
)

var (
	cpuUsage   = expvar.NewFloat("cpu_usage_percent")
	memUsage   = expvar.NewFloat("memory_usage_mb")
	goroutines = expvar.NewInt("goroutines")
)

func init() {
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			memUsage.Set(float64(m.Alloc) / 1024 / 1024)
			goroutines.Set(int64(runtime.NumGoroutine()))

			// Đo CPU (ước lượng)
			cpuUsage.Set(estimateCPUUsage())

			time.Sleep(1 * time.Second)
		}
	}()
}

func estimateCPUUsage() float64 {
	// Đo chính xác hơn bằng runtime/pprof, nhưng cách nhanh:
	var prev, now runtime.MemStats
	runtime.ReadMemStats(&prev)
	time.Sleep(100 * time.Millisecond)
	runtime.ReadMemStats(&now)
	return float64(now.TotalAlloc-prev.TotalAlloc) / 1e8 // heuristic
}
func main() {
	go func() {
		log.Println("pprof đang chạy tại http://localhost:6060/debug/pprof/")
		log.Println("Truy cập: http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof server chết: %v", err) // bắt lỗi ngay
		}
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

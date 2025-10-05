package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
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

func main4() {
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
func main1() {
	go func() {
		fmt.Println("pprof running at :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	r := chi.NewRouter()

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, chi!"))
	})

	http.ListenAndServe(":3000", r)
}

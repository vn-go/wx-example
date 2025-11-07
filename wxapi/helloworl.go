package main

import (
	"fmt"
	"time"

	"github.com/vn-go/wx"
)

type HelloWorld struct {
	PublishOn time.Time
}

func (controller *HelloWorld) New() {
	controller.PublishOn = time.Now()
}

// hello world method POST
// Default method is POST
// any func of HelloWorld has wx.Handler argument that means the handler of the request
// return any type means the response of the request
func (controller *HelloWorld) Hello(h wx.Handler) any {
	return fmt.Sprintf("Hello, World! I was born on %s", controller.PublishOn.Format("2006-01-02 15:04:05"))
}

// This is example of how to use the wx.Handler in the controller
func (controller *HelloWorld) HelloGet(h struct {
	wx.Handler `route:"hello method:GET"`
}) string {
	return fmt.Sprintf("Hello, World! I was born on %s", controller.PublishOn.Format("2006-01-02 15:04:05"))
}

// this is example of hello world with argument , method POST
func (controller *HelloWorld) HelloArg(h wx.Handler, name string) any {
	return fmt.Sprintf("Hello,%s! I was born on %s", name, controller.PublishOn.Format("2006-01-02 15:04:05"))
}

// this is example of hello world with argument in uri path, method GET
func (controller *HelloWorld) HelloUriPath(h struct {
	wx.Handler `route:"hello/{name} method:GET"`
	Name       string
}, name string) any {
	return fmt.Sprintf("Hello,%s! I was born on %s", h.Name, controller.PublishOn.Format("2006-01-02 15:04:05"))
}
func mainHello() {
	//add controller
	err := wx.Routes("/api",

		&HelloWorld{},
	)
	if err != nil {
		panic(err)
	}
	// create server
	server := wx.NewHttpServer("/api", "8080", "localhost")
	// init swagger
	swagger := wx.CreateSwagger(server, "/docs")
	// Show authenication login in swagger
	swagger.OAuth2Password("/api/auth/login")
	// Build swagger
	swagger.Build()
	// start server
	server.Start()
}

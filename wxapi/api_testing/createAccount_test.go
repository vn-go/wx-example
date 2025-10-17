package apitesting

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"wxapi/controllers"

	"runtime/pprof"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/bx"
	"github.com/vn-go/wx"
)

var loginHandler, _ = wx.MakeHandlerFromMethod[controllers.Auth]("Login")

func DoLogin(username, password string) (token string, err error) {

	if err != nil {
		return
	}

	req, err := bx.OnceCall[controllers.Auth]("get-login-request", func() (*http.Request, error) {
		return wx.Mock.FormRequest(loginHandler.GetHttpMethod(), loginHandler.GetUriHandler(), struct {
			Username string
			Password string
		}{
			Username: "admin",
			Password: "/\\dmin123451212",
		})
	})
	if err != nil {
		return
	}
	res := wx.Mock.NewRes()
	loginHandler.Handler().ServeHTTP(res, req)
	if res.Code != 200 {
		err = fmt.Errorf("%s login fail", username)
		return
	}
	data := map[string]string{}
	json.Unmarshal(res.Body.Bytes(), &data)
	token = data["token_type"] + " " + data["access_token"]
	return

}
func TestLogin(t *testing.T) {
	token, err := DoLogin("admin", "/\\dmin123451212")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
func BenchmarkLoginCore(t *testing.B) {
	//go tool pprof BenchmarkLoginCore.prof
	f, err := os.Create("BenchmarkLoginCore.prof")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	// Bắt đầu ghi CPU profile
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	auth := &controllers.Auth{}
	t.Run("login-parallel", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {

			for p.Next() {
				auth.LoginCore(context.Background(), "admin", "/\\dmin123451212")

			}
		})
	})
	t.Run("login", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			auth.LoginCore(context.Background(), "admin", "/\\dmin123451212")
		}
	})
}
func BenchmarkLogin(t *testing.B) {
	//go tool pprof -http=:8080 BenchmarkLogin.prof
	//go test -bench=BenchmarkLogin -benchmem

	//echo top10 | go tool pprof BenchmarkLogin.prof > cpu.txt

	f, err := os.Create("BenchmarkLogin.prof")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	// Bắt đầu ghi CPU profile
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	t.Run("login-parallel", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {

			for p.Next() {
				DoLogin("admin", "/\\dmin123451212")
			}
		})
	})
	t.Run("login", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			DoLogin("admin", "/\\dmin123451212")
		}
	})

}

type userRequest struct {
	Username     string  `json:"username" check:"range:[5:20]"`
	Password     string  `json:"password" check:"range:[5:20]"`
	IsSupperUser bool    `json:"isSupperUser"`
	RoleId       *string `json:"roleId" check:"range:[36:36]"`
}

func TestCreateUser(t *testing.T) {
	token, _ := DoLogin("admin", "/\\dmin123451212")
	createUserAPI, _ := wx.MakeHandlerFromMethod[controllers.Accounts]("UserCreate")
	req, _ := wx.Mock.JsonRequest(createUserAPI.GetHttpMethod(), createUserAPI.GetUriHandler(), &userRequest{
		Username: "user-test",
		Password: "123456",
	})
	req.Header.Add("authorization", token)
	res := wx.Mock.NewRes()
	createUserAPI.Handler().ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)
}
func BenchmarkCreateUser(b *testing.B) {
	f, err := os.Create("BenchmarkCreateUser.prof")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	// Bắt đầu ghi CPU profile
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	start := 2423
	startNext := 2021
	b.Run("default", func(b *testing.B) {
		v := start
		for i := 0; i < b.N; i++ {
			for j := 0; j < 100; j++ {
				token, _ := DoLogin("admin", "/\\dmin123451212")
				createUserAPI, _ := wx.MakeHandlerFromMethod[controllers.Accounts]("UserCreate")
				username := fmt.Sprintf("user-%04d", v)
				req, _ := wx.Mock.JsonRequest(createUserAPI.GetHttpMethod(), createUserAPI.GetUriHandler(), &userRequest{
					Username: username,
					Password: "123456",
				})
				req.Header.Add("authorization", token)
				res := wx.Mock.NewRes()
				createUserAPI.Handler().ServeHTTP(res, req)
				v++
				startNext = v + 1
			}
			//assert.Equal(b, res.Code, 200)

		}
	})
	b.Run("parallel", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			v := startNext
			for p.Next() {
				for j := 0; j < 100; j++ {
					token, _ := DoLogin("admin", "/\\dmin123451212")
					createUserAPI, _ := wx.MakeHandlerFromMethod[controllers.Accounts]("UserCreate")
					username := fmt.Sprintf("user-p-%04d", v)
					req, _ := wx.Mock.JsonRequest(createUserAPI.GetHttpMethod(), createUserAPI.GetUriHandler(), &userRequest{
						Username: username,
						Password: "123456",
					})
					req.Header.Add("authorization", token)
					res := wx.Mock.NewRes()
					createUserAPI.Handler().ServeHTTP(res, req)
					v++
				}
			}
		})
	})
}

/*
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkCreateUser$ wxapi/api_testing

2025/10/10 11:20:23 loading config ./config.yaml
2025/10/10 11:20:23 loading config ./config.yaml,in 0
new db
goos: windows
goarch: amd64
pkg: wxapi/api_testing
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkCreateUser/default-16         	      16	  72002362 ns/op	   31382 B/op	     361 allocs/op
BenchmarkCreateUser/parallel-16        	     199	   5313828 ns/op	   33212 B/op	     393 allocs/op
PASS
ok  	wxapi/api_testing	4.022s

*/

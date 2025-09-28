package main

import (
	"context"
	"core"
	"core/models"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"wxapi/controllers"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/wx"
)

func TestLogin(t *testing.T) {
	//wx.Options.UsePool = true
	handler, err := wx.MakeHandlerFromMethod[controllers.Auth]("Login")
	assert.NoError(t, err)
	req, err := wx.Mock.FormRequest("post", handler.GetUriHandler(), struct {
		Username string
		Password string
	}{
		Username: "tenant-admin",
		Password: "uvsz%#",
	})
	assert.NoError(t, err)
	res := wx.Mock.NewRes()

	handler.Handler().ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)

}
func BenchmarkLoginCore(b *testing.B) {
	auth := &controllers.Auth{}
	b.ResetTimer()
	wx.Options.UsePool = true
	for i := 0; i < b.N; i++ {
		_, err := auth.LoginCore(context.Background(), "test-00001/root", "123456")
		assert.NoError(b, err)
	}
}
func BenchmarkLogin(b *testing.B) {
	handler, err := wx.MakeHandlerFromMethod[controllers.Auth]("Login")
	assert.NoError(b, err)

	uri := handler.GetUriHandler()
	payload := struct {
		Username string
		Password string
	}{
		Username: "test-00001/root",
		Password: "123456",
	}
	req, err := wx.Mock.FormRequest("post", uri, payload)
	if err != nil {
		b.Fatal(err)
	}

	// Tạo response recorder mới
	res := wx.Mock.NewRes()
	b.ResetTimer()
	b.Run("no pool", func(b *testing.B) {
		wx.Options.UsePool = false
		for i := 0; i < b.N; i++ {
			// Tạo request mới

			// Chạy handler
			handler.Handler().ServeHTTP(res, req)

			// (Tuỳ chọn) kiểm tra kết quả hợp lệ
			if res.Code != http.StatusOK {
				b.Fatalf("unexpected status %d", res.Code)
			}
		}
	})
	b.Run("have pool", func(b *testing.B) {
		wx.Options.UsePool = true
		for i := 0; i < b.N; i++ {
			// Tạo request mới

			// Chạy handler
			handler.Handler().ServeHTTP(res, req)

			// (Tuỳ chọn) kiểm tra kết quả hợp lệ
			if res.Code != http.StatusOK {
				b.Fatalf("unexpected status %d", res.Code)
			}
		}
	})

}
func LoginAndGetAccessToken(username, password string) (token string, err error) {
	handler, err := wx.MakeHandlerFromMethod[controllers.Auth]("Login")
	if err != nil {
		return
	}
	req, err := wx.Mock.FormRequest("post", handler.GetUriHandler(), struct {
		Username string
		Password string
	}{
		Username: username, //"tenant-admin",
		Password: password, //"uvsz%#",
	})
	if err != nil {
		return
	}
	res := wx.Mock.NewRes()

	handler.Handler().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", res.Code)
	}
	resData := map[string]any{}
	err = json.Unmarshal(res.Body.Bytes(), &resData)
	if err != nil {
		return
	}
	token = resData["access_token"].(string)
	return
}
func BenchmarkLoginAndGetAccessToken(t *testing.B) {
	defer core.Services.TenantSvc.CloseAllTenants()
	// t.Run("have sync pool", func(b *testing.B) {
	// 	wx.Options.UsePool = true
	// 	for i := 0; i < b.N; i++ {
	// 		_, err := LoginAndGetAccessToken(fmt.Sprintf("test-%04d/root", i), "123456")
	// 		if err != nil {
	// 			t.Error(err)
	// 		}
	// 	}
	// })
	t.Run("no sync pool", func(b *testing.B) {
		wx.Options.UsePool = true
		for i := 0; i < b.N; i++ {
			//_, err := LoginAndGetAccessToken(fmt.Sprintf("test-%04d/root", i), "123456")
			_, err := LoginAndGetAccessToken("admin", "/\\dmin123451212")
			if err != nil {
				t.Error(err)
			}
		}
	})

}

/*
--- no cache---
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLoginAndGetAccessToken$ wxapi

2025/09/27 13:19:06 loading config ./config.yaml
2025/09/27 13:19:06 loading config ./config.yaml,in 0
new db
goos: windows
goarch: amd64
pkg: wxapi
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLoginAndGetAccessToken/no_sync_pool-16         	      13	 120160569 ns/op	   97065 B/op	    1426 allocs/op

PASS
ok  	wxapi	2.795s
----
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLoginAndGetAccessToken$ wxapi

2025/09/27 13:21:13 loading config ./config.yaml
2025/09/27 13:21:13 loading config ./config.yaml,in 0
new db
goos: windows
goarch: amd64
pkg: wxapi
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLoginAndGetAccessToken/no_sync_pool-16         	     100	 119271018 ns/op	   81896 B/op	    1429 allocs/op
---
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLoginAndGetAccessToken$ wxapi

2025/09/27 13:26:47 loading config ./config.yaml
2025/09/27 13:26:47 loading config ./config.yaml,in 4
new db
goos: windows
goarch: amd64
pkg: wxapi
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLoginAndGetAccessToken/no_sync_pool-16         	     100	  95386098 ns/op	   81724 B/op	    1424 allocs/op
BenchmarkLoginAndGetAccessToken/no_sync_pool-16         	     100	 120686987 ns/op	   82487 B/op	    1431 allocs/op
BenchmarkLoginAndGetAccessToken/no_sync_pool-16         	     100	  82011726 ns/op	   81925 B/op	    1420 allocs/op
BenchmarkLoginAndGetAccessToken/no_sync_pool-16         	     100	 118148224 ns/op	   82393 B/op	    1429 allocs/op
*/
func TestCreateTenant(t *testing.T) {
	token, err := LoginAndGetAccessToken("tenant-admin", "uvsz%#")
	assert.NoError(t, err)
	assert.NotNil(t, token)
	handler, err := wx.MakeHandlerFromMethod[controllers.Tenant]("Create")
	assert.NoError(t, err)
	req, err := wx.Mock.JsonRequest(handler.GetHttpMethod(), handler.GetUriHandler(), struct {
		Name string
	}{
		Name: "test-00001",
	})
	assert.NoError(t, err)
	req.Header.Add("authorization", "Bearer "+token)
	res := wx.Mock.NewRes()
	handler.Handler().ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatalf("error")
	}
}

/*
--- use cache inspect token no signature
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkCreateTenant$ wxapi

2025/09/27 13:15:05 loading config ./config.yaml
2025/09/27 13:15:05 loading config ./config.yaml,in 0
new db
close test-0000
goos: windows
goarch: amd64
pkg: wxapi
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkCreateTenant-16    	close test-0000
close test-0001
close test-0002

	3	 345055033 ns/op	  109000 B/op	    1544 allocs/op

PASS
ok  	wxapi	1.936s
*/
func BenchmarkCreateTenant(t *testing.B) {
	defer core.Services.TenantSvc.CloseAllTenants()
	token, err := LoginAndGetAccessToken("tenant-admin", "uvsz%#")
	assert.NoError(t, err)
	assert.NotNil(t, token)
	handler, err := wx.MakeHandlerFromMethod[controllers.Tenant]("Create")
	assert.NoError(t, err)

	for i := 0; i < t.N; i++ {
		req, err := wx.Mock.JsonRequest(handler.GetHttpMethod(), handler.GetUriHandler(), struct {
			Name string
		}{
			Name: fmt.Sprintf("test-%04d", i),
		})
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("authorization", "Bearer "+token)
		res := wx.Mock.NewRes()
		handler.Handler().ServeHTTP(res, req)
		if res.Code != 200 {
			t.Fatalf("error")
		}
	}

}

func BenchmarkCreateTenantParallel(t *testing.B) {
	token, err := LoginAndGetAccessToken("tenant-admin", "uvsz%#")
	assert.NoError(t, err)
	assert.NotNil(t, token)
	handler, err := wx.MakeHandlerFromMethod[controllers.Tenant]("Create")
	assert.NoError(t, err)
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req, err := wx.Mock.JsonRequest(handler.GetHttpMethod(), handler.GetUriHandler(), struct {
				Name string
			}{
				Name: "test-00001",
			})
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("authorization", "Bearer "+token)
			res := wx.Mock.NewRes()
			handler.Handler().ServeHTTP(res, req)
			if res.Code != 200 {
				t.Fatalf("error")
			}
		}
	})

}
func BenchmarkCreateTenantSequential(t *testing.B) {
	token, err := LoginAndGetAccessToken("tenant-admin", "uvsz%#")
	assert.NoError(t, err)
	assert.NotNil(t, token)

	handler, err := wx.MakeHandlerFromMethod[controllers.Tenant]("Create")
	assert.NoError(t, err)

	// Chạy RunParallel nhưng chỉ có 1 goroutine
	t.SetParallelism(1)
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req, err := wx.Mock.JsonRequest(handler.GetHttpMethod(), handler.GetUriHandler(), struct {
				Name string
			}{
				Name: "test-00001",
			})
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("authorization", "Bearer "+token)
			res := wx.Mock.NewRes()
			handler.Handler().ServeHTTP(res, req)
			if res.Code != 200 {
				t.Fatalf("error")
			}
		}
	})
}
func TestAccount(t *testing.T) {
	token, err := LoginAndGetAccessToken("admin", "/\\dmin123451212")
	assert.NoError(t, err)
	handler, err := wx.MakeHandlerFromMethod[controllers.Accounts]("Me")
	assert.NoError(t, err)
	req, err := wx.Mock.JsonRequest(handler.GetHttpMethod(), handler.GetUriHandler(), nil)
	req.Header.Add("authorization", "Bearer "+token)
	assert.NoError(t, err)
	res := wx.Mock.NewRes()
	handler.Handler().ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)
}
func TestRoleCreate(t *testing.T) {
	token, err := LoginAndGetAccessToken("admin", "/\\dmin123451212")
	assert.NoError(t, err)
	handler, err := wx.MakeHandlerFromMethod[controllers.Accounts]("RoleCreate")
	assert.NoError(t, err)
	req, err := wx.Mock.JsonRequest(handler.GetHttpMethod(), handler.GetUriHandler(), models.Role{
		Code: "X-001",
		Name: "X-001",
	})
	req.Header.Add("authorization", "Bearer "+token)
	assert.NoError(t, err)
	res := wx.Mock.NewRes()
	handler.Handler().ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)
}
func BenchmarkRoleCreate(t *testing.B) {
	token, err := LoginAndGetAccessToken("admin", "/\\dmin123451212")
	assert.NoError(t, err)
	handler, err := wx.MakeHandlerFromMethod[controllers.Accounts]("RoleCreate")
	assert.NoError(t, err)
	for i := 0; i < t.N; i++ {
		req, _ := wx.Mock.JsonRequest(handler.GetHttpMethod(), handler.GetUriHandler(), models.Role{
			Code: fmt.Sprintf("X-%04d", i),
			Name: fmt.Sprintf("Role-%04d", i),
		})
		req.Header.Add("authorization", "Bearer "+token)
		// assert.NoError(t, err)
		res := wx.Mock.NewRes()
		handler.Handler().ServeHTTP(res, req)
		// assert.Equal(t, res.Code, 200)
	}

}

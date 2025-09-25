package main

import (
	"context"
	"core"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"wxapi/controllers"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/wx"
)

func TestLogin(t *testing.T) {
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

}
func BenchmarkLoginCore(b *testing.B) {
	auth := &controllers.Auth{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = auth.LoginCore(context.Background(), "test-00001/root", "123456")
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Tạo request mới
		req, err := wx.Mock.FormRequest("post", uri, payload)
		if err != nil {
			b.Fatal(err)
		}

		// Tạo response recorder mới
		res := wx.Mock.NewRes()

		// Chạy handler
		handler.Handler().ServeHTTP(res, req)

		// (Tuỳ chọn) kiểm tra kết quả hợp lệ
		if res.Code != http.StatusOK {
			b.Fatalf("unexpected status %d", res.Code)
		}
	}
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
	for i := 0; i < t.N; i++ {
		LoginAndGetAccessToken(fmt.Sprintf("test-%04d/root", i), "123456")
	}
}

/*
	Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLoginAndGetAccessToken$ wxapi

2025/09/25 21:00:21 loading config ./config.yaml
2025/09/25 21:00:21 loading config ./config.yaml,in 0
new db
goos: windows
goarch: amd64
pkg: wxapi
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLoginAndGetAccessToken-16    	      14	  71682871 ns/op	   87213 B/op	    1412 allocs/op
PASS
ok  	wxapi	1.287s
--- use memcache
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLoginAndGetAccessToken$ wxapi

2025/09/25 21:01:12 loading config ./config.yaml
2025/09/25 21:01:12 loading config ./config.yaml,in 0
new db
goos: windows
goarch: amd64
pkg: wxapi
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLoginAndGetAccessToken-16    	      16	  75083069 ns/op	   91367 B/op	    1406 allocs/op
PASS
ok  	wxapi	1.417s
---- use redis cache---
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkLoginAndGetAccessToken$ wxapi

2025/09/25 21:04:53 loading config ./config.yaml
2025/09/25 21:04:53 loading config ./config.yaml,in 0
new db
goos: windows
goarch: amd64
pkg: wxapi
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkLoginAndGetAccessToken-16    	     100	 120021075 ns/op	   75035 B/op	    1050 allocs/op
PASS
ok  	wxapi	12.162s
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

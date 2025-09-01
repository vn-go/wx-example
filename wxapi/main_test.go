package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/wx"
)

type Example struct {
}
type User struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (ex *Example) FormPost3(ctx *struct {
	wx.Handler `route:"@/files/{*filePath}"`
	FilePath   string
}, data *wx.Form[User]) (*User, error) {
	ret := data.Data
	return &ret, nil
}
func TestFormPost3(t *testing.T) {
	handler, err := wx.MakeHandlerFromMethod[Example]("FormPost3")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := wx.Mock.FormRequest("POST", handler.GetUriHandler()+"dasda/dsad/sad/das/Test.txt", User{
		Code: "adadas",
		Name: "dsadasdasdadad",
	})
	assert.NoError(t, err)
	res := wx.Mock.NewRes()
	fnHandler.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
}
func BenchmarkFormPost3(t *testing.B) {
	handler, err := wx.MakeHandlerFromMethod[Example]("FormPost3")
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	fnHandler := handler.Handler()
	assert.NotNil(t, fnHandler)

	req, err := wx.Mock.FormRequest("POST", handler.GetUriHandler()+"dasda/dsad/sad/das/Test.txt", User{
		Code: "adadas",
		Name: "dsadasdasdadad",
	})
	assert.NoError(t, err)
	res := wx.Mock.NewRes()
	for i := 0; i < t.N; i++ {
		fnHandler.ServeHTTP(res, req)
	}

	//assert.Equal(t, 200, res.Code)
}

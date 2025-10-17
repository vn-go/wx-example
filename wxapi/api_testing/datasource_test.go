package apitesting

import (
	"testing"
	"wxapi/controllers"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/wx"
)

type datasourceRequest struct {
	Name, Fields, Filter string
}

func TestModelDatasourceUser(t *testing.T) {
	token, _ := DoLogin("admin", "/\\dmin123451212")
	createUserAPI, _ := wx.MakeHandlerFromMethod[controllers.DataSource]("Get")
	req, _ := wx.Mock.JsonRequest(createUserAPI.GetHttpMethod(), createUserAPI.GetUriHandler(), &datasourceRequest{
		Name:   "user",
		Filter: "count(id) and email is not null",
		//Fields: "count(id) Total",
	})
	req.Header.Add("authorization", token)
	res := wx.Mock.NewRes()
	createUserAPI.Handler().ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200)
}

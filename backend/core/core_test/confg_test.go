package coretest

import (
	"core"
	"core/services/acc"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/dx"
)

var yamFilePath = "./../../../backend/api/config.yaml"

func TestLoadConfig(t *testing.T) {

	core.Start(yamFilePath)

	core.Services.ConfigSvc.Get()

}

type CacheDataTest struct {
	Name string
	Code string
}

func TestCache(t *testing.T) {
	core.Start(yamFilePath)

	core.Services.CacheSvc.AddObject(t.Context(), "test", "test", &CacheDataTest{
		Name: "AAA",
		Code: "CCC",
	}, 120)
	test := core.Services.CacheSvc.GetObject(t.Context(), "test", "test", &CacheDataTest{})
	assert.NotEmpty(t, test)
	assert.NotEmpty(t, core.Services.AppSvc)
}
func TestAppService(t *testing.T) {

	dx.Options.ShowSql = true
	core.Start(yamFilePath)
	defer core.Services.Close()
	core.Services.AppSvc.InitData()

	fmt.Println(core.Services.AccSvc)
	//fmt.Println(core.Services.DbSvc)
}
func TestLoginService(t *testing.T) {

	dx.Options.ShowSql = true
	core.Start(yamFilePath)
	defer core.Services.Close()
	core.Services.AppSvc.InitData()
	token, err := core.Services.AccSvc.LoginAndGetJWT(t.Context(), "", "admin", "/\\dmin123451212")
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	//fmt.Println(core.Services.DbSvc)
}
func BenchmarkLoginService(b *testing.B) {
	core.Start(yamFilePath)
	defer core.Services.Close()
	core.Services.AppSvc.InitData()
	b.Run("test", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			token, err := core.Services.AccSvc.LoginAndGetJWT(b.Context(), "", "admin", "/\\dmin123451212")
			if err != nil {
				panic(err)
			}
			assert.Greater(b, len(token), 0)
		}
	})
	b.Run("paralell", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				token, err := core.Services.AccSvc.LoginAndGetJWT(b.Context(), "", "admin", "/\\dmin123451212")
				if err != nil {
					panic(err)
				}
				assert.Greater(b, len(token), 0)
			}
		})
	})

}

// UserInfo đại diện cho mỗi đối tượng người dùng trong mảng "data"
type UserInfo struct {
	UserID    string `json:"userID"`
	UserName  string `json:"userName"`
	UserGroup string `json:"userGroup"`
	Email     string `json:"email"`
	Owner     string `json:"owner"`
	Buid      string `json:"buid"`
}

// ConfigData là struct gốc chứa mảng dữ liệu
type ConfigData struct {
	Data []UserInfo `json:"data"`
}

// ReadConfigFromFile đọc nội dung JSON từ filePath và ánh xạ vào struct ConfigData
func ReadConfigFromFile(filePath string) (*ConfigData, error) {
	// 1. Đọc nội dung file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("không thể đọc file %s: %w", filePath, err)
	}

	// 2. Khởi tạo struct để chứa dữ liệu
	var config ConfigData

	// 3. Giải mã (Unmarshal) JSON vào struct
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("không thể giải mã JSON từ file %s: %w", filePath, err)
	}

	return &config, nil
}
func TestRead(t *testing.T) {
	data, err := ReadConfigFromFile(`D:\code\go\wx-example\wx-example\backend\api\data\data3.json`)
	if err != nil {
		panic(err)
	}

	dx.Options.ShowSql = true
	core.Start(yamFilePath)
	defer core.Services.Close()
	core.Services.AppSvc.InitData()
	for _, u := range data.Data {
		err = core.Services.AccSvc.NewUser(t.Context(), "", &acc.UserInfo{
			Username:    u.UserID,
			Email:       u.Email,
			Password:    u.UserID,
			DisplayName: u.UserName,
		})
		if err != nil {
			if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
				if dbErr.IsDuplicateEntryError() {
					continue
				}
			}
			panic(err)
		}
	}
}

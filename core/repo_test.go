package core

import (
	"core/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/dx"
)

func TestRepo(t *testing.T) {
	cfg, err := Config.Load("config.yaml")
	assert.NoError(t, err)
	db, err := dx.Open(cfg.Database.Driver, cfg.Database.DSN)
	assert.NoError(t, err)
	err = db.WithContext(t.Context()).Delete(&models.User{}, "username=?", "root").Error
	assert.NoError(t, err)
	err = Repo.User(db).CreateDefaultUser(db, t.Context(), "123456")
	assert.NoError(t, err)
	err = db.WithContext(t.Context()).Delete(&models.User{}, "username=?", "root").Error
	assert.NoError(t, err)
}
func BenchmarkRepo(t *testing.B) {
	cfg, err := Config.Load("config.yaml")
	if err != nil {
		t.Fail()
	}
	db, err := dx.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	dbHr, err := db.NewDB("hr")
	if err != nil {
		t.Fail()
	}
	for i := 0; i < t.N; i++ {

		err = dbHr.WithContext(t.Context()).Delete(&models.User{}, "username=?", "root").Error
		if err != nil {
			t.Fail()
		}
		err = Repo.User(dbHr).CreateDefaultUser(dbHr, t.Context(), "123456")
		if err != nil {
			t.Fail()
		}
		err = dbHr.WithContext(t.Context()).Delete(&models.User{}, "username=?", "root").Error
		if err != nil {
			t.Fail()
		}
	}

}

/*
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkRepo$ core

goos: windows
goarch: amd64
pkg: core
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkRepo-16    	      28	  38618293 ns/op	  125310 B/op	     574 allocs/op
PASS
ok  	core	1.935s

Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkRepo$ core

goos: windows
goarch: amd64
pkg: core
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkRepo-16    	    4929	    205748 ns/op	   43126 B/op	     124 allocs/op
PASS
ok  	core	1.299s

*/

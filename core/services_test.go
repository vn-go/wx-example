package core

import (
	"context"
	"core/models"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
)

func TestServices(t *testing.T) {
	Users := (&bx.Wire[servicesTypes]{}).WireThenGet(
		newConfig,
		//NewPasswordService,
		newDB,

		newUserRepoSql,
		newUserServiceSql,
		newCacheServiceImpl,
	)
	t.Log(Users)
}
func TestCreateUser(t *testing.T) {
	err := Services.Db.Delete(&models.User{}, "username=?", "test-new-service")
	assert.NoError(t, err.Error)
	user, _ := dx.NewThenSetDefaultValues(func() (*models.User, error) {
		return &models.User{
			Username:     "test-new-service",
			HashPassword: "admminxxxx",
		}, nil
	})
	Services.Users.CreateUser("hrm", t.Context(), user)
	t.Log(user)
}
func BenchmarkCreateUser(t *testing.B) {

	for i := 0; i < t.N; i++ {
		err := Services.TenantSvc.CreateTenant(t.Context(), "hrm", "HRM")
		assert.NoError(t, err)
		Services.Db.Delete(&models.User{}, "username=?", "test-new-service")
		user, _ := dx.NewThenSetDefaultValues(func() (*models.User, error) {
			return &models.User{
				Username:     "test-new-service",
				HashPassword: "admminxxxx",
			}, nil
		})
		Services.Users.CreateUser("hrm", t.Context(), user)

	}

}
func BenchmarkCreateUserParallel(t *testing.B) {
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {
			Services.TenantSvc.CreateTenant(t.Context(), "hrm", "HRM")
			Services.Db.Delete(&models.User{}, "username=?", "test-new-service")
			user, err := dx.NewThenSetDefaultValues(func() (*models.User, error) {
				return &models.User{
					Username:     "test-new-service",
					HashPassword: "admminxxxx",
				}, nil
			})
			if err != nil {
				assert.NoError(t, err)
				t.Fail()
			}
			err = Services.Users.CreateUser("hrm", t.Context(), user)
			if err != nil {
				assert.NoError(t, err)
				t.Fail()
			}
		}
	})

}
func TestLogin(t *testing.T) {
	for i := 0; i < 3; i++ {
		ret, e := Services.AuthSvc.Login("hrm", t.Context(), "root", "123456")
		assert.NoError(t, e)
		assert.NotEmpty(t, ret)
	}

}
func BenchmarkCreateUser2(t *testing.B) {
	go func() {
		for {
			stats := Services.Db.Stats()
			fmt.Printf("Open=%d InUse=%d Idle=%d WaitCount=%d\n",
				stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
			time.Sleep(1 * time.Second)
		}
	}()
	ctx := t.Context()
	for i := 0; i < t.N; i++ {
		user, _ := dx.NewDTO[models.User]()
		user.Username = fmt.Sprintf("user-test-%d", i)
		user.HashPassword = "123456"
		Services.Users.CreateUser("hrm", ctx, user)
		// assert.NoError(t, e)
		// assert.NotEmpty(t, ret)
	}

}
func BenchmarkLogin(t *testing.B) {

	// debug: count open conns
	go func() {
		for {
			stats := Services.Db.Stats()
			fmt.Printf("Open=%d InUse=%d Idle=%d WaitCount=%d\n",
				stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
			time.Sleep(1 * time.Second)
		}
	}()

	//Services.Db.First(&models.User{}) //<--ro rang la co lay du l
	ctx := t.Context()
	for i := 0; i < t.N; i++ {
		Services.AuthSvc.Login("hrm", ctx, fmt.Sprintf("user-test-%d", i), "123456")
		// ret, e := Services.AuthSvc.Login("hrm", ctx, fmt.Sprintf("user-test-%d", i), "123456")
		// assert.NoError(t, e)
		// assert.NotEmpty(t, ret)
	}

}
func BenchmarkDeleteUser(t *testing.B) {

	// debug: count open conns
	go func() {
		for {
			stats := Services.Db.Stats()
			fmt.Printf("Open=%d InUse=%d Idle=%d WaitCount=%d\n",
				stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
			time.Sleep(1 * time.Second)
		}
	}()

	//Services.Db.First(&models.User{}) //<--ro rang la co lay du l
	ctx := t.Context()
	for i := 0; i < t.N; i++ {
		user := &models.User{}
		if err := Services.Db.First(user, "username=?", fmt.Sprintf("user-test-%d", i)); err != nil {
			Services.Users.DeleteUserByUserId("hrm", ctx, user.UserId)
		}

	}

}
func BenchmarkLoginParalel(t *testing.B) {

	// debug: count open conns
	go func() {
		for {
			stats := Services.Db.Stats()
			log.Printf("Open=%d InUse=%d Idle=%d WaitCount=%d\n",
				stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
			fmt.Printf("Open=%d InUse=%d Idle=%d WaitCount=%d\n",
				stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
			time.Sleep(1 * time.Second)
		}
	}()
	t.RunParallel(func(p *testing.PB) {
		i := 0
		for p.Next() {
			Services.AuthSvc.Login("hrm", t.Context(), fmt.Sprintf("user-test-%d", i), "123456")
			i++
		}

	})
	//Services.Db.First(&models.User{}) //<--ro rang la co lay du l
	// ctx := t.Context()
	// for i := 0; i < t.N; i++ {

	// 	// ret, e := Services.AuthSvc.Login("hrm", ctx, fmt.Sprintf("user-test-%d", i), "123456")
	// 	// assert.NoError(t, e)
	// 	// assert.NotEmpty(t, ret)
	// }

}
func TestEncrypt(t *testing.T) {
	txt, err := Services.Encrypt.Encrypt("Hello")
	assert.NoError(t, err)
	log.Println(txt)
}

type User struct {
	UserId string
}

func TestBroker(t *testing.T) {

	defer Services.Broker.CloseAll()
	count := 0

	Services.Broker.Subcribe("delete/user", &User{}, func(ctx context.Context, msg *bx.MsgItem) error {
		fmt.Println(count)
		count++
		return nil
	})

	for i := 0; i < 10; i++ {
		err := Services.Broker.Publish(t.Context(), "delete/user", &User{
			UserId: "123",
		})
		assert.NoError(t, err)
	}

}
func testRun(data any) error {
	//fmt.Println(data)
	return nil
}
func BenchmarkBroker(t *testing.B) {

	defer Services.Broker.CloseAll()
	count := 0
	Services.Broker.Subcribe("delete/user", &User{}, func(ctx context.Context, msg *bx.MsgItem) error {
		fmt.Print(count)
		count++
		msg.Ack()
		return nil
	})
	for i := 0; i < t.N; i++ {

		user := &models.User{}
		Services.Db.First(user, "username=?", fmt.Sprintf("user-test-%d", i))
		err := Services.Broker.Publish(t.Context(), "delete/user", &User{
			UserId: fmt.Sprintf("user %d", i),
		})
		assert.NoError(t, err)
	}

}
func TestFull(t *testing.T) {
	user, err := dx.NewDTO[models.User]()
	assert.NoError(t, err)
	username := "TestFullUser"
	user.Username = username
	user.HashPassword = username
	err = Services.Users.CreateUser("hrm", t.Context(), user)
	assert.NoError(t, err)
	auth, err := Services.AuthSvc.Login("hrm", t.Context(), username, username)
	assert.NoError(t, err)
	assert.NotNil(t, auth)
	err = Services.Users.DeleteUserByUserId("hrm", t.Context(), user.UserId)
	assert.NoError(t, err)
	gUser, err := Services.Users.GetUserByUserId("hrm", t.Context(), user.UserId)
	assert.NoError(t, err)
	assert.Empty(t, gUser)
	auth, err = Services.AuthSvc.Login("hrm", t.Context(), username, username)
	assert.Error(t, err)
	assert.Empty(t, auth)

}
func BenchmarkFull(t *testing.B) {
	for i := 0; i < t.N; i++ {
		Services.TenantSvc.CreateTenant(t.Context(), "hrm", "HRM")
		user, err := dx.NewDTO[models.User]()
		assert.NoError(t, err)
		username := "TestFullUser"
		user.Username = username
		user.HashPassword = username
		err = Services.Users.CreateUser("hrm", t.Context(), user)
		assert.NoError(t, err)
		auth, err := Services.AuthSvc.Login("hrm", t.Context(), username, username)
		assert.NoError(t, err)
		assert.NotNil(t, auth)
		err = Services.Users.DeleteUserByUserId("hrm", t.Context(), user.UserId)
		assert.NoError(t, err)
		gUser, err := Services.Users.GetUserByUserId("hrm", t.Context(), user.UserId)
		assert.NoError(t, err)
		assert.Empty(t, gUser)
		auth, err = Services.AuthSvc.Login("hrm", t.Context(), username, username)
		assert.Error(t, err)
		assert.Empty(t, auth)
	}

}

/*
-- cacche by inmemory
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkFull$ core

2025/09/23 21:07:20 loading config ./config.yaml
2025/09/23 21:07:20 loading config ./config.yaml,in 0
new db
goos: windows
goarch: amd64
pkg: core
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkFull-16    	      10	 107331760 ns/op	   33330 B/op	     403 allocs/op
PASS
ok  	core	1.936s
--- cache by badger
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkFull$ core

2025/09/23 21:12:41 loading config ./config.yaml
2025/09/23 21:12:41 loading config ./config.yaml,in 6
new db
goos: windows
goarch: amd64
pkg: core
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkFull-16    	      10	 106497640 ns/op	   45376 B/op	     602 allocs/op
PASS
ok  	core	2.006s
--- cache by redis cache
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkFull$ core

2025/09/23 21:18:55 loading config ./config.yaml
2025/09/23 21:18:55 loading config ./config.yaml,in 8
new db
goos: windows
goarch: amd64
pkg: core
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkFull-16    	      10	 104645180 ns/op	   35188 B/op	     437 allocs/op
PASS
ok  	core	1.804s
--- cache by memecache
Running tool: C:\Golang\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkFull$ core

2025/09/23 21:28:45 loading config ./config.yaml
2025/09/23 21:28:45 loading config ./config.yaml,in 0
new db
goos: windows
goarch: amd64
pkg: core
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkFull-16    	       9	 129387122 ns/op	   34423 B/op	     426 allocs/op
PASS
ok  	core	1.939s
*/

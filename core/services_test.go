package core

import (
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

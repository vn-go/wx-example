package core

import (
	"context"
	"core/models"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vn-go/bx"
	"github.com/vn-go/dx"
	"go.uber.org/fx"
)

var mysqlDsn = "root:123456@tcp(127.0.0.1:3306)/hrm"

func TestUserServiceAuto(t *testing.T) {
	s, err := bx.NewService[userServiceSql]()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	s.GetUserByUserId("hrm", t.Context(), "73a37553-b83f-4e7f-8486-145696ffb1a3")
}
func TestUserService(t *testing.T) {
	cfg, err := (&Container{}).Get().config.Get() //  (&Container{}).Get() la tao static giong nhu Uber fx
	assert.NoError(t, err)
	db, err := dx.Open(cfg.Database.Driver, cfg.Database.DSN)
	assert.NoError(t, err)

	assert.NoError(t, err)
	c := (&Container{}).Get()
	c.db.Resolve(func() (*dx.DB, error) {
		return db, nil
	})

	userService, err := c.userService.Get()
	assert.NoError(t, err)
	user, _ := dx.NewThenSetDefaultValues(func() (*models.User, error) {
		return &models.User{
			Username: "x-000002",
		}, nil
	})
	err = userService.CreateUser("hrm", t.Context(), user)
	assert.NoError(t, err)
}
func BenchmarkUserService(t *testing.B) {

	t.ResetTimer()
	t.Run("New ID", func(t *testing.B) {
		userService, _ := bx.NewService[userServiceSql]()
		for i := 0; i < t.N; i++ {

			userService.GetUserByUserId("hrm", t.Context(), "73a37553-b83f-4e7f-8486-145696ffb1a3")
		}
	})
	t.Run("New ID init every itegrate", func(t *testing.B) {

		for i := 0; i < t.N; i++ {
			userService, _ := bx.NewService[userServiceSql]()
			userService.GetUserByUserId("hrm", t.Context(), "73a37553-b83f-4e7f-8486-145696ffb1a3")
		}
	})
	t.Run("bx DI", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			c, _ := (&Container{}).NewAndSetValue(func(c *Container) error { // cho nay la tao moi intance

				return nil

			})

			userService, _ := c.userService.Get()
			//user, err :=
			userService.GetUserByUserId("hrm", t.Context(), "73a37553-b83f-4e7f-8486-145696ffb1a3")

		}
	})
	t.Run("bx DI with static container", func(t *testing.B) {
		c := (&Container{}).Get()

		userService, _ := c.userService.Get()
		t.ResetTimer()

		for i := 0; i < t.N; i++ {

			// user, _ := dx.NewThenSetDefaultValues(func() (*models.User, error) {
			// 	return &models.User{
			// 		Username:     uuid.NewString(),
			// 		HashPassword: "password123",
			// 	}, nil
			// })

			//user, err :=
			userService.GetUserByUserId("hrm", t.Context(), "73a37553-b83f-4e7f-8486-145696ffb1a3")
			// assert.NoError(t, err)
			// assert.NotNil(t, user)
		}
	})
	t.Run("fx DI", func(b *testing.B) {
		var userSvc userService

		app := fx.New(
			fx.NopLogger,
			fx.Provide(
				newConfig,
				//NewPasswordService,
				newDB,

				newUserRepoSql,
				newUserServiceSql,
				newCacheServiceImpl,
			),
			fx.Invoke(func(s userService) {
				userSvc = s
			}),
		)

		ctx := context.Background()
		require.NoError(b, app.Start(ctx))
		defer app.Stop(ctx)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			userSvc.GetUserByUserId("hrm", t.Context(), "73a37553-b83f-4e7f-8486-145696ffb1a3")

		}
	})

}
func BenchmarkUserServiceParallel(t *testing.B) {
	var sink int

	//avrg := int64(0)
	t.ResetTimer()
	t.Run("bx DI version 2", func(t *testing.B) {
		t.RunParallel(func(p *testing.PB) {
			bx.Provider(
				newConfig,
				//NewPasswordService,
				newDB,

				newUserRepoSql,
				newUserServiceSql,
				newCacheServiceImpl,
				func() bx.Cache {
					return bx.Cacher.NewInMemoryCache()
				})
			//userService :=
			t.ResetTimer()
			for p.Next() {
				//userService.GetUserByUserId("73a37553-b83f-4e7f-8486-145696ffb1a3")
				sink = bx.GetService[userService]().SayHello()

			}
		})

	})
	t.Run("fx DI", func(b *testing.B) {
		var userSvc userService

		app := fx.New(
			fx.NopLogger,
			fx.Provide(
				newConfig,
				//NewPasswordService,
				newDB,

				newUserRepoSql,
				newUserServiceSql,
				newCacheServiceImpl,
			),
			fx.Invoke(func(s userService) {
				userSvc = s // cho nay se duoc goi bao nhieu lan trong test nay
			}),
		)
		ctx := context.Background()
		require.NoError(b, app.Start(ctx))
		defer app.Stop(ctx)
		t.ResetTimer()
		b.RunParallel(func(p *testing.PB) {

			for p.Next() {
				//userService.GetUserByUserId("73a37553-b83f-4e7f-8486-145696ffb1a3")
				sink = userSvc.SayHello()
			}
		})

	})
	// t.Run("bx DI", func(t *testing.B) {
	// 	t.RunParallel(func(p *testing.PB) {
	// 		for p.Next() {
	// 			c, _ := (&Container{}).NewAndSetValue(func(c *Container) error { // cho nay la tao moi intance
	// 				c.db.Resolve(func() (*dx.DB, error) {
	// 					return db, nil
	// 				})
	// 				c.ctx.Resolve(func() (context.Context, error) {
	// 					return t.Context(), nil
	// 				})

	// 				return nil
	// 			})

	// 			userService, err := c.userService.Get()
	// 			if err != nil {
	// 				t.Log(err)
	// 				t.Fail()
	// 			} else {
	// 				userService.GetUserByUserId("73a37553-b83f-4e7f-8486-145696ffb1a3")
	// 			}

	// 		}
	// 	})

	// })
	t.Run("bx DI with static container", func(t *testing.B) {
		c := (&Container{}).Get()

		t.RunParallel(func(p *testing.PB) {
			for p.Next() {
				userService, _ := c.userService.Get()
				sink = userService.SayHello()

			}
		})

	})
	t.Log(sink)
}

var cnn *dx.DB
var NewDBXOnce sync.Once

func NewDBX(cfg *configInfo) (*dx.DB, error) {
	var err error
	NewDBXOnce.Do(func() {
		cnn, err = dx.Open(cfg.Database.Driver, cfg.Database.DSN)
	})
	return cnn, err
	//return nil, errors.New("loi") // toi thu gia lap loi o day Uber fx bi panic, qua nguy hiem
}
func TestUberFx(t *testing.T) {

	fx.New(
		// cung cấp các dependency gốc
		fx.Provide(
			newConfig,
			//NewPasswordService,
			newDB,

			newUserRepoSql,
			newUserServiceSql,
			newCacheServiceImpl,
		),
		// khởi động app
		fx.Invoke(func(s userService, db *dx.DB, cfg *configInfo) {
			// lam sao dua
			user, _ := dx.NewThenSetDefaultValues(func() (*models.User, error) {
				return &models.User{
					Username: "x-000002",
				}, nil
			})
			s.CreateUser("hrm", t.Context(), user)
		}),
	)
}

var userId = "6095bbb6-2c01-4a87-b4ab-9e8ac5b5179b"

func TestNewService(t *testing.T) {
	userService := bx.ServiceBuidler[userService](
		newConfig,
		//NewPasswordService,
		newDB,

		newUserRepoSql,
		newUserServiceSql,
		newCacheServiceImpl,
		func() bx.Cache {
			return bx.Cacher.NewInMemoryCache()
		},
	).Resovle(func(u userService) any {
		return u
	})

	user, err := userService.GetUserByUserId("hrm", t.Context(), userId)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
func BenchmarkTestDI(b *testing.B) {
	var userData *models.User
	var err error

	b.Run("fx DI", func(b *testing.B) {
		var userSvc userService

		app := fx.New(
			fx.NopLogger,
			fx.Provide(
				newConfig,
				//NewPasswordService,
				newDB,

				newUserRepoSql,
				newUserServiceSql,
				newCacheServiceImpl,
			),
			fx.Invoke(func(s userService) {
				userSvc = s
			}),
		)

		ctx := context.Background()
		require.NoError(b, app.Start(ctx))
		defer app.Stop(ctx)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			userData, err = userSvc.GetUserByUserId("hrm", b.Context(), userId)
			//userService.SayHello()

		}
	})
	b.Run("bx wire", func(b *testing.B) {
		users := (&bx.Wire[servicesTypes]{}).WireThenGet(
			newConfig,
			//NewPasswordService,
			newDB,

			newUserRepoSql,
			newUserServiceSql,
			newCacheServiceImpl,
		).Users
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userData, err = users.GetUserByUserId("hrm", b.Context(), userId)
		}
	})
	b.Run("bx wire in itergrate test", func(b *testing.B) {
		svc := (&bx.Wire[servicesTypes]{}).WireThenGet(
			newConfig,
			//NewPasswordService,
			newDB,

			newUserRepoSql,
			newUserServiceSql,
			newCacheServiceImpl,
		)
		for i := 0; i < b.N; i++ {
			svc.UserRepo.GetUserByUserId(svc.Db, b.Context(), userId)
		}
	})
	b.Log(userData)
	b.Log(err)
}
func BenchmarkUberFx(b *testing.B) {
	var userSvc userService

	app := fx.New(
		fx.Provide(
			newConfig,
			//NewPasswordService,
			newDB,

			newUserRepoSql,
			newUserServiceSql,
			newCacheServiceImpl,
		),
		fx.Invoke(func(s userService) {
			userSvc = s
		}),
	)

	ctx := context.Background()
	require.NoError(b, app.Start(ctx))
	defer app.Stop(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		//userService.GetUserByUserId("73a37553-b83f-4e7f-8486-145696ffb1a3")
		user, err := userSvc.GetUserByUserId("hrm", b.Context(), "73a37553-b83f-4e7f-8486-145696ffb1a3")
		assert.NoError(b, err)
		assert.NotNil(b, user)
	}
}
func BenchmarkBootstrapDi(b *testing.B) {
	b.Run("dx", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c, _ := (&Container{}).NewAndSetValue(func(c *Container) error {

				return nil
			})
			c.userService.Get()
		}
	})
	b.Run("fx", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			app := fx.New(
				fx.NopLogger,

				fx.Provide(
					newConfig,
					//NewPasswordService,
					newDB,

					newUserRepoSql,
					newUserServiceSql,
					newCacheServiceImpl,
				),
				fx.Invoke(func(s userService) {
					_ = s
				}),
			)
			_ = app.Start(context.Background()) // đảm bảo lifecycle chạy
			_ = app.Stop(context.Background())  // dọn tài nguyên
		}
	})
}

package test

import (
	"core"
	"testing"

	"github.com/vn-go/dx"
)

func TestTamperService(t *testing.T) {
	dx.Options.ShowSql = true
	type user struct {
		Username string
	}
	db, err := core.Services.TenantSvc.GetTenant("")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	users, err := dx.QueryItems[user](db, "user(username),where(username like 'user-%')")
	if err != nil {
		panic(err)
	}
	t.Log(users)
	key, err := core.Services.SecretSvc.GenerateMasterKey()
	if err != nil {
		panic(err)
	}
	t.Log(key)
	r, err := core.Services.AuthSvc.Login("", t.Context(), "user-2423", "123456")
	if err != nil {
		panic(err)
	}
	t.Log(r)
}
func BenchmarkTamperService(b *testing.B) {
	type user struct {
		Username string
	}
	db, err := core.Services.TenantSvc.GetTenant("")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	users, err := dx.QueryItems[user](db, "user(username),where(username like 'user-%')")
	if err != nil {
		panic(err)
	}
	b.ResetTimer()

	// b.Run("GenerateMasterKey", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		_, _ = core.Services.SecretSvc.GenerateMasterKey()
	// 	}
	// })
	core.Services.AuthSvc.LoadUserCache("")
	b.Run("Login only", func(b *testing.B) {
		//warmup
		// for _, u := range users {
		// 	_, err := core.Services.AuthSvc.Login("", b.Context(), u.Username, "123456")
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cUser := users[i%200]
			_, err := core.Services.AuthSvc.Login("", b.Context(), cUser.Username, "123456")
			if err != nil {
				panic(err)
			}
		}
	})

}

// BenchmarkTamperService/Login_only-16         	     135	  66060221 ns/op	   33975 B/op	     448 allocs/op
// BenchmarkTamperService/Login_only-16         	      48	  53297308 ns/op	   17873 B/op	     251 allocs/op

func BenchmarkTamperServiceFull(b *testing.B) {
	type user struct {
		Username string
	}

	// --- Load DB / Users ---
	db, err := core.Services.TenantSvc.GetTenant("")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	users, err := dx.QueryItems[user](db, "user(username),where(username like 'user-%')")
	if err != nil {
		panic(err)
	}

	// -------------------------
	// 1️⃣ Benchmark GenerateMasterKey
	b.Run("GenerateMasterKey", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = core.Services.SecretSvc.GenerateMasterKey()
		}
	})

	// -------------------------
	// 2️⃣ Cold login (cache miss)
	b.Run("Login_cold", func(b *testing.B) {
		core.Services.AuthSvc.LoadUserCache("")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cUser := users[i%len(users)]
			_, err := core.Services.AuthSvc.Login("", b.Context(), cUser.Username, "123456")
			if err != nil {
				panic(err)
			}
		}
	})

	// -------------------------
	// 3️⃣ Warm login (cache hit)
	// Warmup: login từng user 1 lần để preload cache
	for _, u := range users {
		_, err := core.Services.AuthSvc.Login("", b.Context(), u.Username, "123456")
		if err != nil {
			panic(err)
		}
	}

	b.Run("Login_warm", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cUser := users[i%len(users)]
			_, err := core.Services.AuthSvc.Login("", b.Context(), cUser.Username, "123456")
			if err != nil {
				panic(err)
			}
		}
	})
}

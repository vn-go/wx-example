package services

import (
	"context"
	"core/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/dx"
)

var mysqlDsn = "root:123456@tcp(127.0.0.1:3306)/hrm"

func TestCreateUser(t *testing.T) {
	db, err := dx.Open("mysql", mysqlDsn)
	assert.NoError(t, err)
	err = db.WithContext(t.Context()).Delete(&models.User{}, "username=?", "test0001").Error
	assert.NoError(t, err)
	userSvc := NewUserServiceSql(db, t.Context())
	user := &models.User{
		Username:     "test0001",
		HashPassword: "test0001",
	}
	err = userSvc.CreateUser(user)
	assert.NoError(t, err)
	assert.Greater(t, user.Id, uint64(0))
}
func BenchmarkCreateUser(b *testing.B) {
	db, err := dx.Open("mysql", mysqlDsn)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := context.Background() // mô phỏng request context
		userSvc := NewUserServiceSql(db, ctx)

		user := &models.User{
			Username:     fmt.Sprintf("bench_user_%d", i),
			HashPassword: "password123",
		}
		if err := userSvc.CreateUser(user); err != nil {
			if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
				if dbErr.ErrorType == dx.Errors.DUPLICATE {
					continue
				}
			}
			b.Fatal(err)
		}
		if user.Id == 0 {
			b.Fatal("user id not set")
		}
	}
}

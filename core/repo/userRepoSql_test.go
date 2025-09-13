package repo

import (
	"core/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/dx"
)

var mysqlDsn = "root:123456@tcp(127.0.0.1:3306)/hrm"

func TestUserRepo(t *testing.T) {
	db, err := dx.Open("mysql", mysqlDsn)
	assert.NoError(t, err)
	repo := NewUserRepoSql(db, t.Context())
	err = db.WithContext(t.Context()).Delete(&models.User{}, "username=?", "root").Error
	assert.NoError(t, err)
	err = repo.CreateDefaultUser("123456")
	assert.NoError(t, err)
	err = db.WithContext(t.Context()).Delete(&models.User{}, "username=?", "root").Error
	assert.NoError(t, err)
}
func BenchmarkUserRepoCreateDefaultUser(b *testing.B) {
	db, err := dx.Open("mysql", mysqlDsn)
	assert.NoError(b, err)
	err = db.Delete(&models.User{}, "username=?", "root").Error
	assert.NoError(b, err)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			repo := NewUserRepoSql(db, b.Context())
			err = db.WithContext(b.Context()).Delete(&models.User{}, "username=?", "root").Error
			assert.NoError(b, err)
			err = repo.CreateDefaultUser("123456")
			assert.NoError(b, err)
			err = db.WithContext(b.Context()).Delete(&models.User{}, "username=?", "root").Error
			assert.NoError(b, err)
		}
	})
}
func TestUserRepoCreateDefaultUserWithHashPass(t *testing.T) {
	db, err := dx.Open("mysql", mysqlDsn)
	assert.NoError(t, err)
	repo := NewUserRepoSql(db, t.Context())
	err = repo.CreateDefaultUser("123456")
	assert.NoError(t, err)

}

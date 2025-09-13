package dxgorm

import (
	"dxgorm/dxmodels"
	_ "dxgorm/dxmodels"
	"dxgorm/gormmodels"
	_ "dxgorm/gormmodels"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	dx "github.com/vn-go/dx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mySqlDsn = "root:123456@tcp(127.0.0.1:3306)/a001"
var mysqlOgrm = "root:123456@tcp(127.0.0.1:3306)/a001?charset=utf8mb4&parseTime=True&loc=Local"
var pgDsn = "postgres://postgres:123456@localhost:5432/a001?sslmode=disable"

func TestDxGetFirstUser(t *testing.T) {
	dxDb, err := dx.Open("mysql", mySqlDsn)
	assert.NoError(t, err)
	assert.NotNil(t, dxDb)
	user, err := dx.NewDTO[dxmodels.User]()
	assert.NoError(t, err)
	err = dxDb.Where("username=?", "admin").First(user)
	assert.NoError(t, err)
}
func BenchmarkDxGetFirstUser(t *testing.B) {
	dxDb, err := dx.Open("mysql", mySqlDsn)
	assert.NoError(t, err)
	assert.NotNil(t, dxDb)
	user, err := dx.NewDTO[dxmodels.User]()
	assert.NoError(t, err)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		err = dxDb.Where("username=?", "admin").First(user)
		assert.NoError(t, err)
	}

}
func TestGormGetFirstUser(t *testing.T) {

	db, err := gorm.Open(mysql.Open(mysqlOgrm), &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	var user gormmodels.User
	err = db.Where("username = ?", "admin").First(&user).Error
	assert.NoError(t, err)

	t.Logf("User: %+v", user)
}
func BenchmarkGormGetFirstUser(t *testing.B) {

	db, err := gorm.Open(mysql.Open(mysqlOgrm), &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	var user gormmodels.User
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		err := db.Where("username = ?", "admin").First(&user).Error
		assert.NoError(t, err)
	}

}
func BenchmarkGetFirstUser(t *testing.B) {
	t.Run("BenchmarkGetFirstUser-dx", func(b *testing.B) {
		BenchmarkDxGetFirstUser(b)
	})
	t.Run("BenchmarkGetFirstUser-gorm", func(b *testing.B) {
		BenchmarkGormGetFirstUser(b)
	})
}
func TestDxFindUser(t *testing.T) {
	dxDb, err := dx.Open("mysql", mySqlDsn)
	assert.NoError(t, err)
	assert.NotNil(t, dxDb)
	user := []dxmodels.User{}
	assert.NoError(t, err)
	for i := 0; i < 5; i++ {
		err = dxDb.Where("username!=?", "admin").Limit(100).Find(&user)
	}

	assert.NoError(t, err)
	fmt.Println(len(user))
}
func TestGormFindUser(t *testing.T) {

	db, err := gorm.Open(mysql.Open(mysqlOgrm), &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	user := []dxmodels.User{}
	err = db.Where("username != ?", "admin").Find(&user).Error
	assert.NoError(t, err)

	t.Logf("User: %+v", user)
}

func BenchmarkDxFindUser(t *testing.B) {
	dxDb, err := dx.Open("mysql", mySqlDsn)
	assert.NoError(t, err)
	assert.NotNil(t, dxDb)

	assert.NoError(t, err)

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		user := []dxmodels.User{}
		dxDb.Where("username!=?", "admin").Find(&user)

	}

}
func BenchmarkGormFindUser(t *testing.B) {

	db, err := gorm.Open(mysql.Open(mysqlOgrm), &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		user := []dxmodels.User{}
		db.Where("username != ?", "admin").Find(&user)
		t.Log(user)

	}

}
func BenchmarkDxUnion(t *testing.B) {
	db, err := dx.Open("mysql", mySqlDsn)
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	// sql := db.Sql(`select tbl.id Id from ( select u.id ID from user u where user.username=?
	// 				union select d.id from department d) tbl`, "admin")

	t.ResetTimer()
	for i := 0; i < t.N; i++ {

		x := []struct {
			Id int
		}{}
		err := db.Sql(`select tbl.id Id from ( 
				select u.id ID from users u where u.username!=?
				union select d.id from departments d) tbl`, "admin").ScanRow(&x)
		assert.NoError(t, err)

	}

}
func BenchmarkGormUnion(t *testing.B) {
	db, err := gorm.Open(mysql.Open(mysqlOgrm), &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// sql := db.Sql(`select tbl.id Id from ( select u.id ID from user u where user.username=?
	// 				union select d.id from department d) tbl`, "admin")
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		x := []struct {
			Id int
		}{}
		err = db.Raw(`select tbl.id Id from ( 
				select u.id ID from users u where u.username!=?
				union select d.id from departments d) tbl`, "admin").Scan(&x).Error
		assert.NoError(t, err)

	}

}
func BenchmarkUnion(b *testing.B) {
	b.Run("dx--union-return-12244-rows", BenchmarkDxUnion)
	b.Run("gorm-union-return-12244-rows", BenchmarkGormUnion)
}
func BenchmarkWithContextSelectUserAndEmailMysqlDx(t *testing.B) {
	db, err := dx.Open("mysql", mySqlDsn)
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	var users = &[]struct {
		Username string
		Email    *string
	}{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		err = db.WithContext(t.Context()).Model(&dxmodels.User{}).Select(
			"concat(username,?,username) username",
			"email", " ",
		).Find(users)
		assert.NoError(t, err)

	}

}
func BenchmarkWithContextSelectUserAndEmailMysqlGorm(t *testing.B) {
	db, err := gorm.Open(mysql.Open(mysqlOgrm), &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	var users = []struct {
		Username string
		Email    *string
	}{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {

		err = db.WithContext(t.Context()).
			Model(&gormmodels.User{}).
			Select("concat(username, ?, username) as username, email", " ").
			Find(&users).Error //12246 dong

		assert.NoError(t, err)
		err = db.WithContext(t.Context()).
			Model(&gormmodels.User{}).
			Select("concat(username, ?, username) as username, email", " ").
			Find(&users).Error //van 12246 dong

		assert.NoError(t, err)
	}

}
func BenchmarkWithContextSelectUserAndEmailMysql(t *testing.B) {
	t.Run("gorm-select-user-email", func(b *testing.B) {
		BenchmarkWithContextSelectUserAndEmailMysqlGorm(b)
	})
	t.Run("dx-select-user-email", func(b *testing.B) {
		BenchmarkWithContextSelectUserAndEmailMysqlDx(b)
	})
}
func TestJoin2ModelDx(t *testing.T) {

	db, err := dx.Open("mysql", mySqlDsn)
	if err != nil {
		t.Fail()
	}
	defer db.Close()
	var users []dxmodels.User

	err = db.Joins("LEFT JOIN department d ON user.Id = department.Id").Limit(10).Find(&users)
	assert.NoError(t, err)

}
func TestJoin2ModelGorm(t *testing.T) {

	db, err := gorm.Open(mysql.Open(mysqlOgrm), &gorm.Config{})
	if err != nil {
		t.Fail()
	}

	var users []gormmodels.User

	err = db.Joins("LEFT JOIN departments  ON users.id = departments.id").Limit(10).Find(&users).Error
	assert.NoError(t, err)

}

func BenchmarkJoin2ModelDxParallel(t *testing.B) {
	db, err := dx.Open("mysql", mySqlDsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()
	//t.ResetTimer()
	t.RunParallel(func(p *testing.PB) {
		for p.Next() {

			var users []dxmodels.User
			err = db.Joins("LEFT JOIN department d ON user.Id = d.Id").Limit(2000).Find(&users)
			assert.NoError(t, err)
		}

	})

}

func BenchmarkJoin2ModelDx(b *testing.B) {
	db, err := dx.Open("mysql", mySqlDsn)
	if err != nil {
		b.Fatalf("failed to open dx db: %v", err)
	}
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []dxmodels.User
		err = db.Joins("LEFT JOIN department ON user.Id = department.Id").
			Limit(2000).
			Find(&users)
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}
	}
}

func BenchmarkJoin2ModelGorm(b *testing.B) {
	db, err := gorm.Open(mysql.Open(mysqlOgrm), &gorm.Config{})
	if err != nil {
		b.Fatalf("failed to open gorm db: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []gormmodels.User
		err = db.Joins("LEFT JOIN departments ON users.id = departments.id").
			Limit(2000).
			Find(&users).Error
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}
	}
}

func BenchmarkJoin2Model(b *testing.B) {
	b.Run("Join2Model-dx-2000-rows", BenchmarkJoin2ModelDx)
	b.Run("Join2Model-gorm-2000-rows", BenchmarkJoin2ModelGorm)
}
func runDxJoin2Model(db *dx.DB) error {

	var users []dxmodels.User
	return db.Joins("LEFT JOIN department ON user.Id = department.Id").
		Limit(2000).
		Find(&users)
}
func runGormJoin2Model(db *gorm.DB) error {

	var users []gormmodels.User
	err := db.Joins("LEFT JOIN departments ON users.id = departments.id").
		Limit(2000).
		Find(&users).Error
	return err
}
func BenchmarkJoin2ModelPal(b *testing.B) {
	dbDx, err := dx.Open("mysql", mySqlDsn)
	if err != nil {
		b.Fail()
	}
	dbGorm, err := gorm.Open(mysql.Open(mysqlOgrm), &gorm.Config{})
	assert.NoError(b, err)
	if err != nil {

		b.Fail()
	}
	defer dbDx.Close()
	b.Run("Join2Model-dx-2000-rows", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				runDxJoin2Model(dbDx)
			}
		})
	})

	b.Run("Join2Model-gorm-2000-rows", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				runGormJoin2Model(dbGorm)
			}
		})
	})
}

package gormmodels

import (
	"time"
)

// BaseModel sử dụng gorm.Model để tự động thêm các trường quan trọng.
// Nó cũng có thể được tùy chỉnh bằng các tag.
type BaseModel struct {
	// GORM sẽ tự động thêm các trường ID, CreatedAt, UpdatedAt, DeletedAt
	// dựa trên quy ước của nó. Vì bạn có các trường tùy chỉnh, chúng ta sẽ định nghĩa lại.
	RecordID    string     `gorm:"type:char(36);uniqueIndex;default:uuid()"`
	CreatedAt   time.Time  `gorm:"index;autoCreateTime"`
	UpdatedAt   *time.Time `gorm:"index;autoUpdateTime"`
	Description *string    `gorm:"size:255"`
}

// Tên bảng mặc định của GORM sẽ là "users" (pluralize).
// Tên trường mặc định của GORM sẽ là "snake_case".
type User struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	UserID       *string `gorm:"type:char(36);uniqueIndex"`
	Email        *string `gorm:"uniqueIndex:uq_email;size:150"`
	Phone        *string `gorm:"size:20"`
	Username     string  `gorm:"size:250;uniqueIndex"`
	HashPassword *string `gorm:"size:100"`

	BaseModel // Nhúng BaseModel vào User
}
type Department struct {
	ID       int    `db:"pk;auto"`
	Name     string `db:"size:100;uk:uq_dept_name"`
	Code     string `db:"size:50;uk:uq_dept_code"`
	ParentID *int
	BaseModel
}

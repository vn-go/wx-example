package models

import "time"

type baseModel struct {
	Id string `db:"pk;size(36); default:uuid()" json:"id"`
	baseInfo
}
type baseCodeName struct {
	Code string `db:"size(50);uk" json:"code"`
	Name string `db:"size(250);idx" json:"name"`
}
type baseInfo struct {
	CreatedOn   time.Time  `db:"idx" json:"createdOn"`
	CreatedBy   string     `db:"size(50);idx" json:"createdBy"`
	ModifiedOn  *time.Time `db:"idx" json:"modifiedOn"`
	ModifiedBy  *string    `db:"size(50);idx" json" json:"modifiedBy"`
	Description string     `db:"size(250);idx" json:"description"`
}

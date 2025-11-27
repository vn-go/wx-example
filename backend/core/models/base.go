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
	CreatedOn   time.Time  `db:"idx" json:"createdOn" private:"true"`
	CreatedBy   string     `db:"size(50);idx" json:"createdBy"  private:"true"`
	ModifiedOn  *time.Time `db:"idx" json:"modifiedOn"  private:"true"`
	ModifiedBy  *string    `db:"size(50);idx" json:"modifiedBy" private:"true"`
	Description string     `db:"size(250);idx" json:"description" `
}

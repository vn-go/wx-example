package models

type CasbinRule struct {
	Id    uint64  `db:"pk;auto"`
	Ptype string  `db:"size:10;ix"` // "p" hoặc "g"
	V0    *string `db:"size:100;ix"`
	V1    *string `db:"size:100;ix"`
	V2    *string `db:"size:100;ix"`
	V3    *string `db:"size:100"`
	V4    *string `db:"size:100"`
	V5    *string `db:"size:100"`
}

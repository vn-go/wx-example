package core

type UserClaims struct {
	Username    string
	UserId      string
	ClaimId     uint64
	RoleId      uint64
	Tenant      string
	IsUpperUser bool
}

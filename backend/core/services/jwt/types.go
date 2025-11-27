package jwt

import jwtV5 "github.com/golang-jwt/jwt/v5"

type JWTClaims[TPayload any] struct {
	jwtV5.RegisteredClaims // standard claims
	Data                   TPayload
}
type Indentifier struct {
	RoleId     *string `json:"roleId"`
	Email      string  `json:"email"`
	Tenant     string  `json:"tenant"`
	UserId     string  `json:"userId"`
	Username   string  `json:"username"`
	ViewPath   string  `json:"-"`
	IsSysAdmin bool    `json:"omitempty"`
}
type IndentifierClaims = JWTClaims[Indentifier]

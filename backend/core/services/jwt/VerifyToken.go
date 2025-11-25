package jwt

import (
	"encoding/base64"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

/*
this function  verify the token with the secret
*/
func (jwtSvc *JwtService) VerifyToken(tokenString string, secret string) (bool, error) {
	// Parse the token
	key, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		// Lỗi này xảy ra nếu chuỗi secret không phải Base64 URL Safe hợp lệ
		return false, errors.New("invalid secret format: cannot decode Base64 URL Safe")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Return the secret
		return []byte(key), nil
	})
	if err != nil {
		return false, err
	}
	// Check if the token is valid
	if token.Valid {
		return true, nil
	}
	return false, errors.New("invalid token")
}

package jwt

import (
	"strings"

	jwtV5 "github.com/golang-jwt/jwt/v5"
)

func (jwtSvc *JwtService) DecodeJwtTokenWithouSignature(token string) (*IndentifierClaims, error) {
	// 1. Initialize the claims structure
	claims := &IndentifierClaims{}

	// 2. Use jwtV5.ParseWithClaims
	// We pass nil as the KeyFunc because we don't want to perform signature validation.
	// The library will correctly parse the claims even without a KeyFunc, but will return
	// a specific error indicating that the token's signature wasn't validated.
	parsedToken, err := jwtV5.ParseWithClaims(token, claims, nil)

	// 3. Check for specific errors that signify a successful *claim extraction* but *failed signature check*

	// The library returns this specific error when no KeyFunc is provided, or the KeyFunc
	// returns an error, OR the signature verification fails. When we pass nil for KeyFunc,
	// this error is guaranteed. We treat this as a *success* for claim extraction.
	if err != nil && !strings.Contains(err.Error(), jwtV5.ErrTokenUnverifiable.Error()) {
		// If the error is anything *other* than the expected "token is unverifiable" error,
		// it means the token was structurally invalid (e.g., wrong number of segments, bad JSON).
		return nil, err
	}

	// 4. Final check: Ensure the parsed token and claims are valid before returning
	if parsedToken != nil && parsedToken.Claims != nil {
		// The claims object is now populated
		return claims, nil
	}

	// Should theoretically not be reached if err is nil or ErrTokenUnverifiable,
	// but serves as a general safeguard.
	return nil, jwtV5.ErrTokenMalformed
}

package core

// import (
// 	"context"
// 	"fmt"
// 	"reflect"
// 	"sort"
// 	"strings"

// 	//"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// type dataJWTService struct {
// 	tenant tenantService
// }
// type DataJWTClaims struct {
// 	jwt.RegisteredClaims
// 	Data       []any    `json:"data"`
// 	DataType   string   `json:"dataType"`
// 	DataFields []string `json:"dataFields"`
// }

// func (s *dataJWTService) ValidateToken(
// 	ctx context.Context,
// 	tokenString string,
// 	tenant string,
// ) (*DataJWTClaims, error) {

// 	db, err := s.tenant.GetTenant(tenant)
// 	if err != nil {
// 		return nil, err
// 	}
// 	secretKey, err := s.tenant.GetSecretKey(ctx, db.DbName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	token, err := jwt.ParseWithClaims(tokenString, &DataJWTClaims{}, func(t *jwt.Token) (interface{}, error) {
// 		// chỉ chấp nhận HS256
// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("invalid signing method")
// 		}
// 		return []byte(secretKey), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(*DataJWTClaims)
// 	if !ok || !token.Valid {
// 		return nil, fmt.Errorf("invalid token claims")
// 	}

// 	return claims, nil
// }

// func (s *dataJWTService) GenerateToken(
// 	ctx context.Context,
// 	user *UserClaims,
// 	fields []string,
// 	data any,
// ) (string, error) {

// 	db, err := s.tenant.GetTenant(user.Tenant)
// 	if err != nil {
// 		return "", err
// 	}
// 	secretKey, err := s.tenant.GetSecretKey(ctx, db.DbName)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Sort fields để đảm bảo payload ổn định & chống reorder attack
// 	sort.Strings(fields)

// 	typ := reflect.TypeOf(data)
// 	val := reflect.ValueOf(data)

// 	if typ.Kind() == reflect.Ptr {
// 		typ = typ.Elem()
// 		val = val.Elem()
// 	}

// 	payload := make([]any, len(fields))

// 	for i, f := range fields {

// 		// Tìm field bằng case-insensitive
// 		field, ok := typ.FieldByNameFunc(func(fn string) bool {
// 			return strings.EqualFold(fn, f)
// 		})
// 		if !ok {
// 			payload[i] = nil
// 			continue
// 		}

// 		fv := val.FieldByIndex(field.Index)

// 		// pointer → handle nil
// 		if fv.Kind() == reflect.Ptr {
// 			if fv.IsNil() {
// 				payload[i] = nil
// 			} else {
// 				payload[i] = fv.Elem().Interface()
// 			}
// 		} else {
// 			payload[i] = fv.Interface()
// 		}
// 	}

// 	claims := DataJWTClaims{
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			Issuer: "eorm-auth",
// 			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
// 			// IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		},
// 		Data:       payload,
// 		DataFields: fields,
// 		DataType:   typ.String(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(secretKey))
// }

// func newDataJWTService(tenant tenantService) *dataJWTService {
// 	return &dataJWTService{
// 		tenant: tenant,
// 	}
// }

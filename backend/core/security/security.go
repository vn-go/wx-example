package security

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type DataContract[T any, TKey comparable] struct {
	Data   T      `json:"data"`
	Key    TKey   `json:"-"`
	Status string `json:"status"`
	Token  string `json:"token"`
}

// CreateToken: serialize data, tạo JWT
func CreateToken[T any](data T, status, secretKey string) (string, error) {
	// Serialize data thành JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"status": status,
		"data":   string(jsonData),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// SignData: gắn token vào DataContract
func SignData[T any, TKey comparable](dc *DataContract[T, TKey], secretKey string) error {
	if dc == nil {
		return ErrNilDataContract
	}
	if secretKey == "" {
		return ErrEmptySecretKey
	}

	token, err := CreateToken(dc.Data, dc.Status, secretKey)
	if err != nil {
		return err
	}

	dc.Token = token
	return nil
}

// NewDataContract: tạo mới DataContract
func NewDataContract[T any, TKey comparable](data T, key TKey) *DataContract[T, TKey] {
	return &DataContract[T, TKey]{
		Data: data,
		Key:  key,
	}
}

// Custom errors
var (
	ErrNilDataContract = errors.New("data contract is nil")
	ErrEmptySecretKey  = errors.New("secret key is empty")
)

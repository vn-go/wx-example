package core

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"

	"github.com/golang-jwt/jwt/v5"
)

type EditClaims[T any, TKey any] struct {
	Data  T      `json:"data"`
	Key   TKey   `json:"-"`
	Token string `json:"token"`
}

type dataSignService struct {
	tenantSvc tenantService
}
type SignedJWTClaims struct {
	jwt.RegisteredClaims
	Data any `json:"data"`
}

func Decrypt(encryptedHex, secretKey string) (string, error) {
	// 1. Decode chuỗi hex thành bytes
	ciphertext, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	// 2. Khởi tạo thuật toán AES và GCM
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 3. Trích xuất Nonce
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext quá ngắn")
	}
	// Nonce được đính kèm ở đầu chuỗi ciphertext
	nonce, encryptedText := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 4. Giải mã dữ liệu (Decrypt) và Xác thực (Authenticate)
	// Nếu dữ liệu bị giả mạo, hàm Open sẽ trả về lỗi
	plaintextBytes, err := aesGCM.Open(nil, nonce, encryptedText, nil)
	if err != nil {
		// Lỗi thường là "cipher: message authentication failed" nếu dữ liệu bị thay đổi
		return "", fmt.Errorf("giải mã hoặc xác thực thất bại: %v", err)
	}

	// Trả về chuỗi plaintext
	return string(plaintextBytes), nil
}
func Encrypt(plaintext, secretKey string) (string, error) {
	// 1. Khởi tạo thuật toán AES
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// 2. Sử dụng chế độ GCM (Galois/Counter Mode)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 3. Tạo Nonce (Số chỉ được sử dụng một lần)
	// Nonce phải là duy nhất cho mỗi lần mã hóa với cùng một khóa
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 4. Mã hóa dữ liệu (Encrypt)
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// Trả về chuỗi mã hóa dưới dạng chuỗi hex để dễ dàng truyền tải qua JSON/HTTP
	return hex.EncodeToString(ciphertext), nil
}

/*
	 this function will make a new HMAC-SHA256 with Key data

	 Example:
	 	data= EditClaims[models.users,struct {
			username string
			userid string
			...
		}] {
			Data: User,
			Key: struct {
				username: User.username,
				....
			}
		}
*/
func (s *dataSignService) SignData(ctx context.Context, tenant string, data any) error {

	secretKey, err := s.tenantSvc.GetSecretKey(ctx, tenant)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	keyField, ok := typ.FieldByName("Key")
	if !ok {
		return fmt.Errorf("key field of %T was not found", data)
	}
	dataVal := val.FieldByName("Data")
	if !dataVal.IsValid() {
		return fmt.Errorf("argument Data field of %T was not found", data)
	}
	keyType := keyField.Type
	keyVal := reflect.New(keyType).Elem()
	for i := 0; i < keyType.NumField(); i++ {

		fieldOfValInData := dataVal.FieldByName(keyType.Field(i).Name)
		if fieldOfValInData.IsValid() {
			keyVal.Field(i).Set(fieldOfValInData)
		}
	}
	claims := SignedJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Issuer: "eorm-auth",
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			// IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Data: keyVal.Interface(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//return token.SignedString([]byte(secretKey))
	// // mak
	// bff, err := json.Marshal(keyValue)
	// if err != nil {
	// 	return err
	// }

	// // Trả về chuỗi băm dưới dạng hex
	// token, err := Encrypt(string(bff), secretKey)
	// if err != nil {
	// 	return err
	// }
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}
	val.FieldByName("Token").SetString(accessToken)
	return nil

}

func newDataSignService(tenantSvc tenantService) *dataSignService {
	return &dataSignService{
		tenantSvc: tenantSvc,
	}
}

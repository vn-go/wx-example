package core

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

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

type getKeyFieldResult struct {
	keyFieldsIndex            [][]int
	keyFieldIndex             []int
	keyFieldIndexInMasterData [][]int
	TokenFieldIndex           []int
	dataFieldIndex            []int
	keyType                   reflect.Type
}

func (s *dataSignService) getKeyFieldNoCache(data any) (*getKeyFieldResult, error) {
	ret := &getKeyFieldResult{
		keyFieldsIndex:            [][]int{},
		keyFieldIndexInMasterData: [][]int{},
	}
	typ := reflect.TypeOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()

	}

	typeOfDataField, ok := typ.FieldByName("Data")
	if !ok {
		return nil, fmt.Errorf("field Data was not found in %T", data)
	}
	ret.dataFieldIndex = typeOfDataField.Index
	keyField, ok := typ.FieldByName("Key")
	ret.keyFieldIndex = keyField.Index
	if !ok {
		return nil, fmt.Errorf("key field of %T was not found", data)
	}
	tokenField, ok := typ.FieldByName("Token")
	if !ok {
		return nil, fmt.Errorf("Token field of %T was not found", data)
	}
	ret.TokenFieldIndex = tokenField.Index
	ret.keyType = keyField.Type
	for i := 0; i < keyField.Type.NumField(); i++ {
		ret.keyFieldsIndex = append(ret.keyFieldsIndex, keyField.Type.Field(i).Index)
		fieldName := keyField.Type.Field(i).Name

		dataField, ok := typeOfDataField.Type.FieldByName(fieldName)
		if !ok {
			return nil, fmt.Errorf("%s was not found in %T", keyField.Type.Field(i).Name, data)
		}
		ret.keyFieldIndexInMasterData = append(ret.keyFieldIndexInMasterData, append(typeOfDataField.Index, dataField.Index...))
	}

	return ret, nil
}

type initGetKeyField struct {
	val  *getKeyFieldResult
	err  error
	once sync.Once
}

var initGetKeyFieldCache sync.Map

func (s *dataSignService) getKeyField(data any) (*getKeyFieldResult, error) {
	typ := reflect.TypeOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	a, _ := initGetKeyFieldCache.LoadOrStore(typ, &initGetKeyField{})
	i := a.(*initGetKeyField)
	i.once.Do(func() {
		i.val, i.err = s.getKeyFieldNoCache(data)
	})
	return i.val, i.err

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
func (s *dataSignService) SignData(ctx context.Context, user *UserClaims, data any) error {

	secretKey := user.UserId
	// if err != nil {
	// 	return err
	// }
	val := reflect.ValueOf(data)
	//typ := reflect.TypeOf(data)
	if val.Kind() == reflect.Ptr {
		//typ = typ.Elem()
		val = val.Elem()
	}

	fiedlInfo, err := s.getKeyField(data)
	if err != nil {
		return err
	}
	keyVal := reflect.New(fiedlInfo.keyType).Elem()
	if err != nil {
		return err
	}
	for i, f := range fiedlInfo.keyFieldsIndex {
		fieldOfValInData := val.FieldByIndex(fiedlInfo.keyFieldIndexInMasterData[i]) //dataVal.FieldByIndex(fiedlInfo.dataFeidlIndex[i])
		if fieldOfValInData.IsValid() {
			keyVal.FieldByIndex(f).Set(fieldOfValInData)
		}
	}
	// for i := 0; i < keyType.NumField(); i++ {

	// 	fieldOfValInData := dataVal.FieldByName(keyType.Field(i).Name)
	// 	if fieldOfValInData.IsValid() {
	// 		keyVal.Field(i).Set(fieldOfValInData)
	// 	}
	// }
	claims := SignedJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		Data:             keyVal.Interface(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}
	val.FieldByName("Token").SetString(accessToken)
	return nil

}
func (s *dataSignService) Verify(ctx context.Context, user *UserClaims, data any) error {
	secret := user.UserId
	// if err != nil {
	// 	return err
	// }
	keyFieldInfo, err := s.getKeyField(data)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	tokenFieldValue := val.FieldByIndex(keyFieldInfo.TokenFieldIndex)
	if !tokenFieldValue.IsValid() {
		return fmt.Errorf("can not read token from %T", data)
	}
	tokenString := tokenFieldValue.String()

	// Parse token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid token claims")
	}

	rawData, err := json.Marshal(claims["data"])
	if err != nil {
		return err
	}

	dataVal := val.FieldByIndex(keyFieldInfo.keyFieldIndex)
	if !dataVal.IsValid() {
		return fmt.Errorf("invalid data in %T", data)
	}

	ptr := dataVal.Addr() // *T

	// ⬇ Unmarshal trực tiếp vào field
	if err := json.Unmarshal(rawData, ptr.Interface()); err != nil {
		return err
	}

	for i, f := range keyFieldInfo.keyFieldsIndex {

		field := val.FieldByIndex(keyFieldInfo.keyFieldIndexInMasterData[i])

		if field.IsValid() {
			valField := dataVal.FieldByIndex(f)
			if valField.IsValid() {
				data := valField.Interface()
				fmt.Println(data)
				field.Set(valField)
			}

		}
	}

	return nil
}

func newDataSignService(tenantSvc tenantService) *dataSignService {
	return &dataSignService{
		tenantSvc: tenantSvc,
	}
}

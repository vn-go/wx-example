package data

import (
	"context"
	"core/services/jwt"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	jwtV5 "github.com/golang-jwt/jwt/v5"
)

type getKeyFieldResult struct {
	keyFieldsIndex            [][]int
	keyFieldIndex             []int
	keyFieldIndexInMasterData [][]int
	TokenFieldIndex           []int
	dataFieldIndex            []int
	keyType                   reflect.Type
}
type DataSignService struct {
}
type SignedJWTClaims struct {
	jwtV5.RegisteredClaims
	Data any `json:"data"`
}

func (s *DataSignService) getKeyFieldNoCache(data any) (*getKeyFieldResult, error) {
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

func (s *DataSignService) getKeyField(data any) (*getKeyFieldResult, error) {
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
func (s *DataSignService) SignData(ctx context.Context, user *jwt.Indentifier, data any) error {

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
		RegisteredClaims: jwtV5.RegisteredClaims{},
		Data:             keyVal.Interface(),
	}

	token := jwtV5.NewWithClaims(jwtV5.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return err
	}
	val.FieldByName("Token").SetString(accessToken)
	return nil

}
func (s *DataSignService) Verify(ctx context.Context, user *jwt.Indentifier, data any) error {
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
	token, err := jwtV5.Parse(tokenString, func(t *jwtV5.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtV5.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwtV5.MapClaims)
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

func NewDataSignService() *DataSignService {
	return &DataSignService{}
}

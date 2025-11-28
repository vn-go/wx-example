package data

import (
	"context"
	"core/services/errs"
	"core/services/jwt"
	"fmt"
	"reflect"
	"sync"

	json "github.com/json-iterator/go"

	jwtV5 "github.com/golang-jwt/jwt/v5"
)

type getKeyFieldResult struct {
	keyFieldsIndex            [][]int
	keyFieldIndex             []int
	keyFieldIndexInMasterData [][]int
	TokenFieldIndex           []int
	dataFieldIndex            []int
	keyType                   reflect.Type
	StatuFieldIndex           []int
}
type DataSignService struct {
	errs.ErrService
}
type SignedJWTClaims struct {
	jwtV5.RegisteredClaims
	Data   any    `json:"data"`
	Status string `json:"status"`
}

func (s *DataSignService) getKeyFieldNoCache(typ reflect.Type) (*getKeyFieldResult, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	ret := &getKeyFieldResult{
		keyFieldsIndex:            [][]int{},
		keyFieldIndexInMasterData: [][]int{},
	}
	// typ := reflect.TypeOf(data)
	// if typ.Kind() == reflect.Ptr {
	// 	typ = typ.Elem()

	// }

	typeOfDataField, ok := typ.FieldByName("Data")
	if !ok {
		return nil, fmt.Errorf("field Data was not found in %s", typ.String())
	}
	ret.dataFieldIndex = typeOfDataField.Index
	Status, ok := typ.FieldByName("Status")

	if !ok {
		return nil, fmt.Errorf("Status field of %s was not found", typ.String())
	}
	ret.StatuFieldIndex = Status.Index
	keyField, ok := typ.FieldByName("Key")

	if !ok {
		return nil, fmt.Errorf("key field of %s was not found", typ.String())
	}
	ret.keyFieldIndex = keyField.Index
	tokenField, ok := typ.FieldByName("Token")
	if !ok {
		return nil, fmt.Errorf("Token field of %s was not found", typ)
	}
	ret.TokenFieldIndex = tokenField.Index
	ret.keyType = keyField.Type
	for i := 0; i < keyField.Type.NumField(); i++ {
		ret.keyFieldsIndex = append(ret.keyFieldsIndex, keyField.Type.Field(i).Index)
		fieldName := keyField.Type.Field(i).Name

		// fieldIndex, err := s.getFieldIndexByFieldName(typeOfDataField, fieldName)
		// if err != nil {
		// 	return nil, err
		// }
		// fmt.Println(fieldIndex)
		dataField, ok := typeOfDataField.Type.FieldByName(fieldName)
		if !ok {
			return nil, fmt.Errorf("%s was not found in %s", keyField.Type.Field(i).Name, typ)
		}

		ret.keyFieldIndexInMasterData = append(ret.keyFieldIndexInMasterData, append(typeOfDataField.Index, dataField.Index...))
	}

	return ret, nil
}

// func (s *DataSignService) getFieldIndexByFieldName(typeOfDataField reflect.StructField, fieldName string) ([]int, error) {
// 	ret, ok := typeOfDataField.Type.FieldByNameFunc(func(name string) bool {
// 		return strings.EqualFold(name, fieldName)
// 	})
// 	if ok {
// 		return ret.Index, nil
// 	} else {
// 		for i := 0; i < typeOfDataField.Type.NumField(); i++ {
// 			fmt.Println(typeOfDataField.Type.Field(i).Name)
// 			if typeOfDataField.Type.Field(i).Anonymous {
// 				innerField, err := s.getFieldIndexByFieldName(typeOfDataField.Type.Field(i), fieldName)
// 				if err != nil {
// 					return nil, err
// 				}
// 				return append(typeOfDataField.Index, innerField...), nil
// 			}
// 		}
// 	}
// 	return nil, fmt.Errorf("%s was not found in %T", fieldName, typeOfDataField.Type)
// }

type initGetKeyField struct {
	val  *getKeyFieldResult
	err  error
	once sync.Once
}

var initGetKeyFieldCache sync.Map

func (s *DataSignService) getKeyField(typ reflect.Type) (*getKeyFieldResult, error) {

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	a, _ := initGetKeyFieldCache.LoadOrStore(typ, &initGetKeyField{})
	i := a.(*initGetKeyField)
	i.once.Do(func() {
		i.val, i.err = s.getKeyFieldNoCache(typ)
	})
	if i.err != nil || i.val == nil {
		initGetKeyFieldCache.Delete(typ)
		return nil, i.err
	}
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

	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)
	if val.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
		if typ.Kind() != reflect.Struct {

			panic("data must be a struct")
		}
	}

	fiedlInfo, err := s.getKeyField(typ)
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
	statusVal := val.FieldByIndex(fiedlInfo.StatuFieldIndex).String()
	claims := SignedJWTClaims{
		RegisteredClaims: jwtV5.RegisteredClaims{},
		Data:             keyVal.Interface(),
		Status:           statusVal,
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
	typ := reflect.TypeOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()

	}
	keyFieldInfo, err := s.getKeyField(typ)
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
			return nil, s.BadRequest("invalid token")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return s.BadRequest("invalid token")
	}

	claims, ok := token.Claims.(jwtV5.MapClaims)
	if !ok || !token.Valid {
		return s.BadRequest("invalid token") //errors.New("invalid token claims")
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

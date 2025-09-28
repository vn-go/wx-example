package core

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/vn-go/dx"
)

type errType struct {
}
type ActionError struct {
	ServiceName string      `json:"service_name"`
	Action      string      `json:"action"`
	Code        string      `json:"code"`
	DbErr       *dx.DbError `json:"_"`
	Err         error       `json:"_"`
}

func (a *ActionError) Error() string {
	v, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("%s,%s,%s", a.ServiceName, a.Action, a.Code)
	}
	return string(v)
}
func (e *errType) Create(service any, funcName string, err error) error {
	typ := reflect.TypeOf(service)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if dbErr := dx.Errors.IsDbError(err); dbErr != nil {
		return &ActionError{
			ServiceName: typ.Name(),
			Action:      funcName,
			Code:        dx.String(dbErr.ErrorType),
			DbErr:       dbErr,
		}
	}
	return &ActionError{
		ServiceName: typ.Name(),
		Action:      funcName,
		Err:         err,
	}
}

var Errors = &errType{}

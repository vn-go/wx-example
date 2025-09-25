package be

import (
	"fmt"
	"reflect"
	"testing"
)

func LoadData(data any) {
	// data phải là **T
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		panic("LoadData requires a non-nil pointer to a pointer")
	}

	// trỏ tới *T
	elem := v.Elem()
	if elem.IsNil() {
		// lấy type T
		typ := elem.Type().Elem()
		// tạo *T mới
		newVal := reflect.New(typ)
		// gán vào *T (elem)
		elem.Set(newVal)
	}
}
func NewData() (ret *struct{ X string }) {
	LoadData(&ret)
	return ret
}
func TestXxx(t *testing.T) {
	r := NewData()
	fmt.Println(r)
}

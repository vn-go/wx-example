package internal

import (
	"errors"
	"reflect"
	"sync"
)

type oneCallItem[T any] struct {
	Val T
	Err error
}

type initOnceCall[T any] struct {
	Once sync.Once
	Item oneCallItem[T] // struct trực tiếp, không dùng *pointer
}

var cacheOnceCall sync.Map // map[reflect.Type]*initOnceCall[T]

func OnceCall[TService any, TResult any](fnName string, fn func() (TResult, error)) (TResult, error) {
	serviceType := reflect.TypeFor[TService]()
	if serviceType.Kind() != reflect.Struct {
		var zero TResult
		return zero, errors.New("serviceInstance arg of OnceCall must be struct or pointer to struct")
	}

	// key là reflect.Type, không cần build string -> tránh alloc
	actual, _ := cacheOnceCall.LoadOrStore(serviceType, &initOnceCall[TResult]{})
	onceCall := actual.(*initOnceCall[TResult])

	onceCall.Once.Do(func() {
		onceCall.Item.Val, onceCall.Item.Err = fn()
	})

	return onceCall.Item.Val, onceCall.Item.Err
}

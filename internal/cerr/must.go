package must

import (
	"fmt"
	"reflect"
)

func BeNotEmpty(value string, field string) {
	if value == "" {
		panic(fmt.Errorf("must not be empty: %q", field))
	}
}

func BeGreater(value int, gt int, field string) {
	if value <= gt {
		panic(fmt.Errorf("must be greater %d: %q", gt, field))
	}
}

func BeNotNil(ptr any, field string) {
	// this does not work for interfaces or other special types
	errf := "must not be nil: %q"
	if ptr == nil {
		panic(fmt.Errorf(errf, field))
	}

	// this code should only run on application start up, therefore performance is not an issue
	value := reflect.ValueOf(ptr)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		if value.IsNil() {
			panic(fmt.Errorf(errf, field))
		}
	default:
		return
	}
}

package scheduler

import (
	"fmt"
	"reflect"
)

func mustFieldOffset[T any](field string) uintptr {
	var a T
	f, ok := reflect.TypeOf(a).FieldByName(field)
	if !ok {
		panic(fmt.Errorf("unable to get %s field of %T", field, a))
	}
	return f.Offset
}

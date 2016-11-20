package util

import (
	//"fmt"
	"reflect"
)

func ToMap(any interface{}) map[string]string {
	vt := reflect.TypeOf(any)
	vv := reflect.ValueOf(any)

	ret := make(map[string]string)

	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)

		switch reflect.TypeOf(f).Kind() {
		case reflect.Slice, reflect.Array:

		case reflect.Map:

		}

		ret[f.Name] = vv.FieldByName(f.Name).String()
	}
	return ret
}

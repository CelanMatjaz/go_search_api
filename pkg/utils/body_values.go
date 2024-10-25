package utils

import (
	"fmt"
	"reflect"
)

func GetValuesFromBody(body any, prepend []any) []any {
	v := reflect.ValueOf(body)

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				continue
			} else {
				field = field.Elem()
			}
		}

		if field.Kind() == reflect.Struct {
			prepend = GetValuesFromBody(field.Interface(), prepend)
			continue
		}

		if bodyTag, ok := fieldType.Tag.Lookup("body"); !ok || bodyTag == "omit" {
			continue
		}

		prepend = append(prepend, getValues(field)...)
	}

	return prepend
}

func getValues(val reflect.Value) []any {
	values := make([]any, 0)

	kind := val.Kind()
	switch kind {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		values = append(values, val.Int())
		break

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		values = append(values, val.Uint())
		break

	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			values = append(values, val.Index(i).Int())
		}
		break

	case reflect.String:
		values = append(values, val.String())
		break

	default:
		panic(fmt.Sprintf("case for reflect.Kind = %d is not provided", kind))
	}

	return values
}

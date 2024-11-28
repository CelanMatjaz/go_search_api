package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func GetValuesFromBody(body any, neededTagValue string, prepend []any) []any {
	v := reflect.ValueOf(body)

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// TODO: create a more generic solution
		if fieldType.Tag.Get("db") == "account_id" {
			continue
		}

		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				continue
			} else {
				field = field.Elem()
			}
		}

		if field.Kind() == reflect.Struct {
			prepend = GetValuesFromBody(field.Interface(), neededTagValue, prepend)
			continue
		}

		if bodyTag, ok := fieldType.Tag.Lookup("body"); !ok || bodyTag == "omit" || !strings.Contains(bodyTag, neededTagValue) {
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		values = append(values, val.Int())
		break

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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

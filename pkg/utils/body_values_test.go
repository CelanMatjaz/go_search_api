package utils

import (
	"reflect"
	"testing"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

func TestGetValuesFromBody(t *testing.T) {
	t.Parallel()

	type Struct1 struct {
		value1 int    `body:"create"`
		value2 []int  `body:"update,create"`
		value3 string `body:"update"`
	}

	type Struct2 struct {
		Struct1
		value4 int    `body:""`
		value5 []int  `body:"create"`
		value6 string `body:"update"`
	}

	type Struct3 struct {
		Struct2
		*types.WithTags
	}

	createDefaultStruct := func() Struct3 {
		return Struct3{
			Struct2: Struct2{
				Struct1: Struct1{
					value1: 1,
					value2: []int{2, 3, 4},
					value3: "value 5",
				},
				value4: 6,
				value5: []int{7, 8, 9},
				value6: "value 10",
			},
			WithTags: &types.WithTags{
				TagIds: []int{11, 12, 13, 14},
			},
		}
	}

	createAndCheckSlices := func(tagValue string, expected []any) (bool, []any, []any) {
		values1 := GetValuesFromBody(createDefaultStruct(), tagValue, []any{})
		if !reflect.DeepEqual(values1, expected) {
			return false, values1, expected
		}

		prepend := []any{1, 2, 3}
		values2 := GetValuesFromBody(createDefaultStruct(), tagValue, prepend)
		prepend = append(prepend, expected...)
		if !reflect.DeepEqual(values2, prepend) {
			return false, values2, prepend
		}

		return true, nil, nil
	}

	convertSliceIntegers := func(slice []any) []any {
		for index, value := range slice {
			t := reflect.TypeOf(value)
			v := reflect.ValueOf(value)
			if kind := t.Kind(); kind >= 2 && kind <= 11 {
				slice[index] = v.Int()
			}
		}

		return slice
	}

	t.Run("with empty body", func(t *testing.T) {
		expected := convertSliceIntegers([]any{1, 2, 3, 4, "value 5", 6, 7, 8, 9, "value 10", 11, 12, 13, 14})
		if ok, generated, expected := createAndCheckSlices("", expected); !ok {
			t.Fatalf("expected and generated values are not equal\nexpected:  %v\ngenerated: %v", expected, generated)
		}
	})

	t.Run("with select body", func(t *testing.T) {
		expected := convertSliceIntegers([]any{1, 2, 3, 4, 7, 8, 9, 11, 12, 13, 14})
		if ok, generated, expected := createAndCheckSlices("create", expected); !ok {
			t.Fatalf("expected and generated values are not equal\nexpected:  %v\ngenerated: %v", expected, generated)
		}
	})

	t.Run("with create body", func(t *testing.T) {
		expected := convertSliceIntegers([]any{1, 2, 3, 4, 7, 8, 9, 11, 12, 13, 14})
		if ok, generated, expected := createAndCheckSlices("create", expected); !ok {
			t.Fatalf("expected and generated values are not equal\nexpected:  %v\ngenerated: %v", expected, generated)
		}
	})

	t.Run("with update body", func(t *testing.T) {
		expected := convertSliceIntegers([]any{2, 3, 4, "value 5", "value 10"})
		if ok, generated, expected := createAndCheckSlices("update", expected); !ok {
			t.Fatalf("expected and generated values are not equal\nexpected:  %v\ngenerated: %v", expected, generated)
		}
	})

	type TestCase struct {
		expectedValues [][]any
	}

	t.Run("tag", func(t *testing.T) {
		tagValues := []string{"", "create", "update"}
		testCase := TestCase{
			expectedValues: [][]any{
				{0, "label", "color1"},
				{"label", "color1"},
				{"label", "color1"},
			},
		}

		tag := types.CreateTag(123, "label", "color1")
		for i, expectedValues := range testCase.expectedValues {
			expected := convertSliceIntegers(expectedValues)
			generated := GetValuesFromBody(tag, tagValues[i], []any{})
			if !reflect.DeepEqual(generated, expected) {
				t.Fatalf("expected and generated values are not equal for tag value \"%v\"\nexpected:  %v\ngenerated: %v", tagValues[i], expected, generated)
			}
		}
	})
}

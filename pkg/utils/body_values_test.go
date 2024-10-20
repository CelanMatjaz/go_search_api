package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetValuesFromBody(t *testing.T) {
	type TestCase struct {
		structWithValues any
		expectedValues   []any
		shouldFail       bool
	}

	type Struct1 struct {
		value1 int    `body:""`
		value2 []int  `body:""`
		value3 string `body:""`
	}

	type Struct2 struct {
		Struct1
		value4 int    `body:""`
		value5 []int  `body:""`
		value6 string `body:""`
	}

	testCases := []TestCase{
		{
			structWithValues: Struct1{value1: 1, value2: []int{2, 3, 4}, value3: "test"},
			expectedValues:   []any{int64(1), int64(2), int64(3), int64(4), "test"},
			shouldFail:       false,
		},
		{
			structWithValues: Struct1{value1: 1, value2: []int{2, 3, 4}, value3: ""},
			expectedValues:   []any{1, 2, 3, 4, "test"},
			shouldFail:       true,
		},
		{
			structWithValues: Struct2{
				Struct1: Struct1{value1: 1, value2: []int{2, 3, 4}, value3: "test"},
				value4:  5,
				value5:  []int{6, 7, 8},
				value6:  "test",
			},
			expectedValues: []any{int64(1), int64(2), int64(3), int64(4), "test", int64(5), int64(6), int64(7), int64(8), "test"},
			shouldFail:     false,
		},
		{
			structWithValues: Struct2{
				Struct1: Struct1{value1: 1, value2: []int{2, 3, 4}, value3: "test"},
				value4:  5,
				value5:  []int{6, 7, 8},
				value6:  "test",
			},
			expectedValues: []any{1, 2, 3, 4, "test"},
			shouldFail:     true,
		},
	}

	printSlices := func(expectedValues []any, generatedValues []any) {
		fmt.Printf("expected values  %v\n", expectedValues)
		fmt.Printf("generated values %v\n", generatedValues)
	}

	compareSlices := func(expectedValues []any, generatedValues []any) error {
		if len(expectedValues) != len(generatedValues) {
			return fmt.Errorf("slices are not of same length")
		}

		for i := 0; i < len(expectedValues); i++ {
			expectedType := reflect.TypeOf(expectedValues[i])
			generatedType := reflect.TypeOf(generatedValues[i])

			if expectedType != generatedType {
				return fmt.Errorf("expected and generated values are not of equal type at index %d ", i)
			}

			if !reflect.DeepEqual(expectedValues[i], generatedValues[i]) {
				return fmt.Errorf("expected and generated values are not of equal in value at index %d ", i)
			}
		}

		return nil
	}

	for i, tc := range testCases {
		values := GetValuesFromBody(tc.structWithValues)
		if err := compareSlices(tc.expectedValues, values); (err != nil) && !tc.shouldFail {
			printSlices(tc.expectedValues, values)
			t.Errorf(
				"GetValuesFromBody() failed for test case with index %d, should fail: %t\nerror: %s",
				i, tc.shouldFail, err.Error())
			break
		}
	}
}

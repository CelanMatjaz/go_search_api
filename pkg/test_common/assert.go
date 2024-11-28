package testcommon

import (
	"reflect"
	"testing"
)

func Assert(t *testing.T, statement bool, message string, args ...any) {
	if !statement {
		t.Fatalf(message, args...)
	}
}

func AssertNotError(t *testing.T, err error, message string) {
	if err != nil {
		if len(message) > 0 {
			t.Fatalf("unexpected error: %s, %s", message, err.Error())
		} else {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	}
}

func AssertError(t *testing.T, error error, expectedError error) {
	if error == nil {
		t.Fatalf("expected error but got nil")
	}

	if expectedError != nil && error.Error() != expectedError.Error() {
		t.Fatalf("unexpected error\nexpected: %s (%s)\nactual:   %s\n",
			expectedError.Error(),
			reflect.TypeOf(expectedError).Name(),
			error.Error())
	}
}

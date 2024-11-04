package testcommon

import (
	"fmt"
	"testing"
)

func Assert(t *testing.T, statement bool, message string, args ...any) {
	if !statement {
		t.Fatalf(message, args...)
	}
}

func AssertError(t *testing.T, err error, message string) {
	if err != nil {
		t.Fatalf(fmt.Sprintf("%s, %%s", message), err.Error())
	}
}

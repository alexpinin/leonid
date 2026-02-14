package testutil

import (
	"reflect"
	"testing"
)

func Equal(t testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Mismatched values. Expected: %v Actual: %v", expected, actual)
	}
}

package testutil

import (
	"reflect"
	"slices"
	"testing"
)

func Equal(t testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Mismatched values. Expected: %v Actual: %v", expected, actual)
	}
}

func EqualSlice[S ~[]E, E comparable](t testing.TB, expected S, actual S) {
	if !slices.Equal(expected, actual) {
		t.Fatalf("Mismatched values. Expected: %v Actual: %v", expected, actual)
	}
}

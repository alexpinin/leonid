package testutil

import (
	"errors"
	"maps"
	"reflect"
	"slices"
	"strings"
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

func EqualMap[M1, M2 ~map[K]V, K, V comparable](t testing.TB, expected M1, actual M2) {
	if !maps.Equal(expected, actual) {
		t.Fatalf("Mismatched values. Expected: %v Actual: %v", expected, actual)
	}
}

func NotEqual(t testing.TB, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected values to be different. Expected: %v Actual: %v", expected, actual)
	}
}

func Contains[T comparable](t testing.TB, s []T, v T) {
	for _, sv := range s {
		if sv == v {
			return
		}
	}
	t.Fatalf("Expected slice to contain value. Expected: %v Actual: %v", v, s)
}

func HasPrefix(t testing.TB, prefix, s string) {
	if !strings.HasPrefix(s, prefix) {
		t.Fatalf("Expected string to start with prefix. Expected prefix: %v Actual string: %v", prefix, s)
	}
}

func HasSuffix(t testing.TB, suffix, s string) {
	if !strings.HasSuffix(s, suffix) {
		t.Fatalf("Expected string to end with suffix. Expected suffix: %v Actual string: %v", suffix, s)
	}
}

func ErrorIs(t testing.TB, expected, actual error) {
	if !errors.Is(actual, expected) {
		t.Fatalf("Mismatched errors. Expected: %v Actual: %v", expected, actual)
	}
}

func ErrorContains(t testing.TB, expected string, err error) {
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Fatalf("Expected error to contain substring. Expected: %v Actual: %v", expected, err)
	}
}

func MustSucceed(t testing.TB, err error) {
	if err != nil {
		t.Fatalf("Expected success. Actual: %v", err)
	}
}

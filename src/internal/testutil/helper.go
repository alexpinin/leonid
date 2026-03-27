package testutil

import (
	"encoding/json"
	"testing"
)

func MustMarshalJson(t *testing.T, v any) []byte {
	j, err := json.Marshal(v)
	MustSucceed(t, err)
	return j
}

func MustUnmarshalJson[T any](t *testing.T, v []byte) T {
	var o T
	err := json.Unmarshal(v, &o)
	MustSucceed(t, err)
	return o
}

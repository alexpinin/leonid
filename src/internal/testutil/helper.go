package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

func MustCreateRequest(t *testing.T, method, path string, data any) *http.Request {
	j := MustMarshalJson(t, data)
	body := bytes.NewReader(j)
	req, err := http.NewRequest(method, path, body)
	MustSucceed(t, err)
	return req
}

func MustReadBody[T any](t *testing.T, body io.ReadCloser) T {
	inst := new(T)
	all, err := io.ReadAll(body)
	err = json.Unmarshal(all, &inst)
	MustSucceed(t, err)
	return *inst
}

func LongString(length int) string {
	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		sb.WriteRune('a')
	}
	return sb.String()
}

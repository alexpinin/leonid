package handler

import (
	"testing"

	"leonid/src/internal/testutil"
)

func TestBasicHandlerNextHandle(t *testing.T) {
	t.Run("should call next handler if it's not nil", func(t *testing.T) {
		next := &mockHandler{}
		sut := &basicHandler{}
		sut.setNext(next)

		_ = sut.nextHandle(nil, nil, nil, &UpdateState{})

		testutil.Equal(t, 1, next.handleCount)
	})
	t.Run("should not call next and fail if it's nil", func(t *testing.T) {
		sut := &basicHandler{}
		_ = sut.nextHandle(nil, nil, nil, &UpdateState{})
	})
}

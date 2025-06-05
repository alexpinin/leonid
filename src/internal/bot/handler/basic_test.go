package handler

import (
	"leonid/testutil"
	"testing"
)

func Test_basicHandler_nextHandle(t *testing.T) {
	t.Run("should call next handler if it's not nil", func(t *testing.T) {
		runLog := make([]string, 0)
		sut := &basicHandler{}
		sut.setNext(&mockHandler{runLog: &runLog})

		sut.nextHandle(nil, nil, &UpdateContext{})

		expectedRunLog := []string{"handle: " + testUpdateToStr(&UpdateContext{})}
		testutil.Equal(t, expectedRunLog, runLog)
	})
	t.Run("should not call next and fail if it's nil", func(t *testing.T) {
		sut := &basicHandler{}
		sut.nextHandle(nil, nil, &UpdateContext{})
	})
}

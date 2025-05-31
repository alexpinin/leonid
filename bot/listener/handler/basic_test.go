package handler

import (
	"leonid/testutil"
	"testing"
)

func Test_basicHandler_handleNext(t *testing.T) {

	t.Run("should call next handler if it's not nil", func(t *testing.T) {
		m := &mockHandler{}
		h := &basicHandler{}
		h.SetNext(m)

		h.handleNext(nil, nil, &UpdateContext{})

		testutil.Equal(t, &UpdateContext{}, m.updateContext)
	})

	//t.Run("should not call next handler if it's nil", func(t *testing.T) {
	//	m := &mockHandler{}
	//	h := &basicHandler{}
	//	h.SetNext(nil)
	//
	//	h.handleNext(nil, nil, &UpdateContext{})
	//
	//	testutil.Equal(t, *UpdateContext(nil), m.UpdateContext)
	//})

}

package listener

import (
	"leonid/testutil"
	"os"
	"reflect"
	"testing"
)

func Test_Handler(t *testing.T) {
	t.Run("it should use a specific handlers order", func(t *testing.T) {
		_ = os.Setenv("DB_FILE", "DB_FILE")
		sut := NewHandler()

		expected := []string{
			"*handler.InputGuard",
			"*handler.ChatChecker",
			"*handler.ChatActivator",
			"*handler.AuthGuard",
			"*handler.CallGuard",
			"*handler.QuotaGuard",
			"*handler.MessageCleaner",
			"*handler.MessageSender",
		}

		i := 0
		next := sut.handlerHead
		for next != nil && i < len(expected) {
			testutil.Equal(t, expected[i], reflect.TypeOf(next).String())
			next = next.GetNext()
			i++
		}
	})
}

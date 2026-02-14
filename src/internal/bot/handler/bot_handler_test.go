package handler

import (
	"os"
	"reflect"
	"testing"

	"leonid/src/internal/bot/service"
	"leonid/testutil"
)

func init() {
	_ = os.Setenv("DB_FILE", "DB_FILE")
}

func Test_Handler(t *testing.T) {
	t.Run("it should use a specific handlers order", func(t *testing.T) {
		sut := NewBotHandler(
			&service.ConfigService{},
			&service.QuotaService{},
			&service.OpenAIMessageService{},
		)

		expected := []string{
			"*handler.inputGuard",
			"*handler.chatChecker",
			"*handler.chatActivator",
			"*handler.authGuard",
			"*handler.callGuard",
			"*handler.quotaGuard",
			"*handler.messageSender",
		}

		i := 0
		next := sut.handlerHead
		for next != nil && i < len(expected) {
			testutil.Equal(t, expected[i], reflect.TypeOf(next).String())
			next = next.getNext()
			i++
		}
	})
}

package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UpdateContext struct {
	*models.Update
	IsInputValid bool
	IsChatActive bool
}

type UpdateHandler interface {
	Handle(context.Context, *bot.Bot, *UpdateContext)
	SetNext(UpdateHandler)
}

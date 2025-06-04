package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UpdateContext struct {
	*models.Update
	IsChatActive bool
	IsPassActive bool
}

type UpdateHandler interface {
	Handle(context.Context, *bot.Bot, *UpdateContext)
	SetNext(UpdateHandler)
	GetNext() UpdateHandler
}

package app

import (
	"context"
	"gobot1/dispatch"

	"github.com/gotd/td/telegram"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Client *telegram.Client
	DM     *dispatch.DMDispatcher
}

func New(parent context.Context, db *gorm.DB, client *telegram.Client) *App {
	return &App{
		DB:     db,
		Client: client,
		DM:     dispatch.NewDMDispatcher(parent),
	}
}

package app

import (
	"github.com/gotd/td/telegram"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Client *telegram.Client
}

func New(db *gorm.DB, client *telegram.Client) *App {
	return &App{
		DB:     db,
		Client: client,
	}
}

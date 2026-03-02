package handlers

import (
	"gobot1/router"

	"github.com/gotd/td/telegram"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Client *telegram.Client
	Router router.EventRouter
}

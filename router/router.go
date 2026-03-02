package router

import (
	"context"

	"github.com/gotd/td/tg"
)

type EventRouter interface {
	RouteNewMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error
}

package handlers

import (
	"context"

	"github.com/gotd/td/tg"
)

func (h *Handler) HandleNewMessage(
	ctx context.Context,
	e tg.Entities,
	u *tg.UpdateNewMessage,
) error {
	return h.Router.RouteNewMessage(ctx, e, u)
}

package app

import (
	"context"
	"fmt"

	"github.com/gotd/td/tg"
)

func (a *App) OnDMNewMessage(ctx context.Context, e tg.Entities, msg *tg.Message, u *tg.UpdateNewMessage) error {
	fmt.Println("[DM]", msg.Message)
	return nil
}

func (a *App) OnGroupNewMessage(ctx context.Context, e tg.Entities, msg *tg.Message, u *tg.UpdateNewMessage) error {
	fmt.Println("[GROUP]", msg.Message)
	return nil
}

func (a *App) OnChannelNewMessage(ctx context.Context, e tg.Entities, msg *tg.Message, u *tg.UpdateNewMessage) error {
	fmt.Println("[CHANNEL]", msg.Message)
	return nil
}

func (a *App) OnUnknownNewMessage(ctx context.Context, e tg.Entities, msg *tg.Message, u *tg.UpdateNewMessage) error {
	fmt.Println("[UNKNOWN]", msg.Message)
	return nil
}

func (a *App) OnNonMessageUpdate(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
	// update.Message был не *tg.Message (например сервисный/другой класс)
	fmt.Printf("[NON-MESSAGE] %T\n", u.Message)
	return nil
}

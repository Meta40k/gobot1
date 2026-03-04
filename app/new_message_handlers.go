package app

import (
	"context"
	"fmt"

	"github.com/gotd/td/tg"
)

func (a *App) OnDMNewMessage(
	ctx context.Context,
	e tg.Entities,
	msg *tg.Message,
	u *tg.UpdateNewMessage,
) error {

	userID, ok := extractUserIDForDM(msg)
	if !ok {
		// Не смогли определить пользователя — игнорируем или логируем
		return nil
	}

	// Кладём задачу в очередь конкретного пользователя.
	a.DM.Enqueue(userID, func(workerCtx context.Context) error {
		return a.HandleDMNewMessage(workerCtx, e, msg, u)
	})

	// Важно: мы НЕ обрабатываем сообщение прямо здесь.
	return nil
}

// HandleDMNewMessage — реальная обработка DM, выполняется в goroutine пользователя.
func (a *App) HandleDMNewMessage(
	ctx context.Context,
	e tg.Entities,
	msg *tg.Message,
	u *tg.UpdateNewMessage,
) error {
	fmt.Println("[DM handled in worker]", msg.Message)
	return nil
}

// extractUserIDForDM достаёт userID для лички.
// Для DM чаще всего достаточно msg.FromID, но fallback на PeerID полезен.
func extractUserIDForDM(msg *tg.Message) (int64, bool) {
	if pu, ok := msg.FromID.(*tg.PeerUser); ok && pu.UserID != 0 {
		return pu.UserID, true
	}
	if pu, ok := msg.PeerID.(*tg.PeerUser); ok && pu.UserID != 0 {
		return pu.UserID, true
	}
	return 0, false
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

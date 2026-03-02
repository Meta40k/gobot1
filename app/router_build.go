package app

import "gobot1/router"

func (a *App) BuildRouter() router.EventRouter {
	mr := router.NewMessageRouter()

	// 3 уровень: назначаем обработчики веток
	mr.OnDM = a.OnDMNewMessage
	mr.OnGroup = a.OnGroupNewMessage
	mr.OnChannel = a.OnChannelNewMessage
	mr.OnUnknown = a.OnUnknownNewMessage
	mr.OnNonMessage = a.OnNonMessageUpdate

	return mr
}

package router

import (
	"context"
	"fmt"

	"github.com/gotd/td/tg"
)

type ChatKind int

const (
	ChatUnknown ChatKind = iota
	ChatDM
	ChatLegacyGroup
	ChatSupergroup
	ChatChannel
)

type NewMessageHandler func(
	ctx context.Context,
	entities tg.Entities,
	msg *tg.Message,
	update *tg.UpdateNewMessage,
) error

type MessageRouter struct {
	OnDM         NewMessageHandler
	OnGroup      NewMessageHandler // legacy group (PeerChat) + supergroup (Megagroup)
	OnChannel    NewMessageHandler // broadcast channel
	OnUnknown    NewMessageHandler
	OnNonMessage func(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error
}

func NewMessageRouter() *MessageRouter {
	fmt.Println("----------Отработал метод func NewMessageRouter() *MessageRouter---------------------")
	return &MessageRouter{}
}

func (r *MessageRouter) RouteNewMessage(
	ctx context.Context,
	entities tg.Entities,
	update *tg.UpdateNewMessage,
) error {
	fmt.Println("---Отработал метод func (r *MessageRouter) RouteNewMessage")
	// update.Message — это MessageClass (может быть не *tg.Message)
	msg, ok := update.Message.(*tg.Message)
	if !ok {
		if r.OnNonMessage != nil {
			return r.OnNonMessage(ctx, entities, update)
		}
		return nil
	}

	kind, err := classifyPeer(msg.PeerID, entities)
	if err != nil {
		// если не смогли классифицировать, пусть решает Unknown handler
		if r.OnUnknown != nil {
			return r.OnUnknown(ctx, entities, msg, update)
		}
		return nil
	}

	switch kind {
	case ChatDM:
		if r.OnDM != nil {
			return r.OnDM(ctx, entities, msg, update)
		}
	case ChatLegacyGroup, ChatSupergroup:
		if r.OnGroup != nil {
			return r.OnGroup(ctx, entities, msg, update)
		}
	case ChatChannel:
		if r.OnChannel != nil {
			return r.OnChannel(ctx, entities, msg, update)
		}
	default:
		if r.OnUnknown != nil {
			return r.OnUnknown(ctx, entities, msg, update)
		}
	}

	return nil
}

// classifyPeer определяет тип чата по PeerID + Entities.
func classifyPeer(peer tg.PeerClass, entities tg.Entities) (ChatKind, error) {
	switch p := peer.(type) {
	case *tg.PeerUser:
		return ChatDM, nil

	case *tg.PeerChat:
		return ChatLegacyGroup, nil

	case *tg.PeerChannel:
		ch, ok := entities.Channels[p.ChannelID]
		if !ok || ch == nil {
			// Entities могут быть неполными, но для роутера это “неизвестно”
			return ChatUnknown, fmt.Errorf("channel %d not found in entities", p.ChannelID)
		}
		// В gotd tg.Channel имеет флаги Megagroup/Broadcast (bool поля в структуре).
		if ch.Broadcast {
			return ChatChannel, nil
		}
		if ch.Megagroup {
			return ChatSupergroup, nil
		}
		// На всякий случай: channel без обоих флагов редок, считаем Unknown.
		return ChatUnknown, nil

	default:
		return ChatUnknown, fmt.Errorf("unsupported peer type %T", peer)
	}
}

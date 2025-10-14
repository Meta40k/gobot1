package util

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func FindChannelByID(ctx context.Context, client *telegram.Client, channelID int64) (*tg.Channel, error) {
	if client == nil {
		return nil, fmt.Errorf("client not initialized")
	}
	api := client.API()

	res, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetDate: 0,
		OffsetID:   0,
		OffsetPeer: &tg.InputPeerEmpty{},
		Limit:      200,
	})
	if err != nil {
		return nil, fmt.Errorf("get dialogs: %w", err)
	}

	switch d := res.(type) {
	case *tg.MessagesDialogs:
		for _, ch := range d.Chats {
			if c, ok := ch.(*tg.Channel); ok && c.ID == channelID {
				return c, nil
			}
		}
	case *tg.MessagesDialogsSlice:
		for _, ch := range d.Chats {
			if c, ok := ch.(*tg.Channel); ok && c.ID == channelID {
				return c, nil
			}
		}
	default:
		return nil, fmt.Errorf("unexpected dialogs type %T", res)
	}

	return nil, fmt.Errorf("channel %d not found in first page of dialogs (нужно состоять в канале)", channelID)
}

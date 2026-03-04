package main

import (
	"context"
	"fmt"
	"gobot1/app"
	"gobot1/handlers"
	"gobot1/router"
	"gobot1/util"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gotd/td/examples"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type program struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	config     *Config
	botClient  *telegram.Client
	handler    *handlers.Handler
	dispatcher tg.UpdateDispatcher
	flow       auth.Flow
	db         *gorm.DB
	router     router.EventRouter
}

type TerminalWithPassword struct {
	examples.Terminal
	PasswordValue string
}

func (t TerminalWithPassword) Password(ctx context.Context) (string, error) {
	return t.PasswordValue, nil
}

func (p *program) Initialize(args []string) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	p.config = cfg

	p.ctx, p.cancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	p.dispatcher = tg.NewUpdateDispatcher()
	p.botClient = telegram.NewClient(p.config.APIID, p.config.APIHash, telegram.Options{
		SessionStorage: &session.FileStorage{Path: "tg.session"},
		UpdateHandler:  p.dispatcher,
	})

	p.flow = auth.NewFlow(
		TerminalWithPassword{
			Terminal:      examples.Terminal{PhoneNumber: p.config.Phone},
			PasswordValue: p.config.CloudPassword, // может быть пустым, если 2FA выключен
		},
		auth.SendCodeOptions{},
	)

	db, err := openDB()
	if err != nil {
		return err
	}
	if err := migrate(db); err != nil {
		return err
	}

	p.db = db
	application := app.New(p.ctx, p.db, p.botClient)

	p.router = application.BuildRouter()

	fmt.Printf("%T\n", p.router)

	p.handler = &handlers.Handler{Client: p.botClient, DB: p.db, Router: p.router}
	p.registerHandlers()

	return nil
}

func (p *program) wait() {
	p.wg.Wait()
}

func (p *program) Run() error {
	if p.config == nil || p.botClient == nil {
		return fmt.Errorf("program not initialized")
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		_ = p.botClient.Run(p.ctx, func(ctx context.Context) error {
			if err := p.botClient.Auth().IfNecessary(ctx, p.flow); err != nil {
				return err
			}

			util.Whoami(ctx, p.botClient)

			//metakChannel, _ := strconv.ParseInt(os.Getenv("METAK"), 10, 64)
			//p.startAdminLogPoller(metakChannel, time.Hour)
			//err := p.pollAdminLogOnce(ctx, metakChannel)
			//if err != nil {
			//	return err
			//}

			<-ctx.Done()
			return ctx.Err()
		})
	}()
	return nil
}

func (p *program) pollAdminLogOnce(ctx context.Context, channelID int64) error {

	channel, err := util.FindChannelByID(ctx, p.botClient, channelID)
	if err != nil {
		return fmt.Errorf("find channel: %w", err)
	}

	InputChannel := &tg.InputChannel{ChannelID: channel.ID, AccessHash: channel.AccessHash}

	api := p.botClient.API()
	AdminLog, err := api.ChannelsGetAdminLog(ctx, &tg.ChannelsGetAdminLogRequest{
		Channel: InputChannel,
		Limit:   50,
	})
	if err != nil {
		return fmt.Errorf("find channel: %w", err)

	}

	fmt.Printf("AdminLog poll: %d events\n", len(AdminLog.Events))

	for i, event := range AdminLog.Events {
		timestamp := time.Unix(int64(event.Date), 0).Format(time.RFC3339)
		fmt.Printf("%2d) %s by user %d action=%T (eventID=%d)\n",
			i+1, timestamp, event.UserID, event.Action, event.ID)
	}

	rows := make([]AdminlogEvent, 0, len(AdminLog.Events))

	for _, event := range AdminLog.Events {
		ts := time.Unix(int64(event.Date), 0).UTC()
		var actor *int64
		if event.UserID != 0 {
			id := event.UserID
			actor = &id
		}

		rows = append(rows, AdminlogEvent{
			ChannelID:   channel.ID,
			EventID:     event.ID,
			EventTsUTC:  ts,
			ActorUserID: actor,
			ActionType:  fmt.Sprintf("%T", event.Action),
		})
	}

	if len(rows) > 0 {
		if err := p.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rows).Error; err != nil {
			return fmt.Errorf("insert adminlog events: %w", err)
		}
	}

	return nil
}

func (p *program) startAdminLogPoller(channelID int64, interval time.Duration) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-p.ctx.Done():
				return
			case <-ticker.C:
				_ = p.pollAdminLogOnce(p.ctx, channelID) // TODO: обработать ошибку/лог
			}
		}
	}()
}

func (p *program) registerHandlers() {
	p.dispatcher.OnNewMessage(p.handler.HandleNewMessage)
}

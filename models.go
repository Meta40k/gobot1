// models.go
package main

import (
	"time"

	"gorm.io/datatypes"
)

// AdminlogEvent — одна запись из admin log.
type AdminlogEvent struct {
	ChannelID   int64          `gorm:"primaryKey;not null"`
	EventID     int64          `gorm:"primaryKey;not null"` // BIGINT, уникален в рамках канала
	EventTsUTC  time.Time      `gorm:"not null;index:idx_channel_time,sort:desc"`
	ActorUserID *int64         `gorm:"index"`      // инициатор (может быть nil)
	ActionType  string         `gorm:"not null"`   // например: "*tg.ChannelAdminLogEventActionParticipantLeave"
	RawAction   datatypes.JSON `gorm:"type:jsonb"` // опционально: сырой payload
}

func (AdminlogEvent) TableName() string { return "tg.adminlog_events" }

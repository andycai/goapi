package adminlog

import (
	"context"

	"github.com/andycai/goapi/events"
	"github.com/andycai/goapi/models"
	"github.com/andycai/goapi/pkg/event"
)

func subscribeEvents(bus *event.EventBus) {
	event.Subscribe(bus, event.EventHandler[events.EventAddOperationLog](commandAddLog))
}

// commandAddLog 写入操作日志命令
func commandAddLog(ctx context.Context, event events.EventAddOperationLog) error {
	log := models.AdminLog{
		UserID:     event.UserID,
		Username:   event.Username,
		Action:     event.Action,
		Resource:   event.Resource,
		ResourceID: event.ResourceID,
		Details:    event.Details,
		IP:         event.IP,
		UserAgent:  event.UserAgent,
		CreatedAt:  app.DB.NowFunc(),
	}

	if err := app.DB.Create(&log).Error; err != nil {
		return err
	}

	return nil
}

package adminlog

import (
	"context"

	"github.com/andycai/goapi/core/event"
	"github.com/andycai/goapi/models"
)

func init() {
	event.Subscribe(app.Bus, event.EventHandler[AddOperationLogEvent](commandAddLog))
}

// commandAddLog 写入操作日志命令
func commandAddLog(ctx context.Context, event AddOperationLogEvent) error {
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

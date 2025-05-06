package adminlog

import (
	"context"
	"fmt"

	"github.com/andycai/goapi/core/event"
	"github.com/andycai/goapi/models"
	"github.com/gofiber/fiber/v2"
)

func init() {
	event.Subscribe(app.Bus, event.EventHandler[AddOperationLogEvent](commandAddLog))
}

func commandAddLog(ctx context.Context, event AddOperationLogEvent) error {
	c := ctx.Value("fiberCtx").(*fiber.Ctx)
	currentUser := app.CurrentUser(c)

	if currentUser.ID == 0 {
		return fmt.Errorf("登录已过期，请重新登录")
	}

	log := models.AdminLog{
		UserID:     currentUser.ID,
		Username:   currentUser.Username,
		Action:     event.Action,
		Resource:   event.Resource,
		ResourceID: event.ResourceID,
		Details:    event.Details,
		IP:         c.IP(),
		UserAgent:  c.Get("User-Agent"),
		CreatedAt:  app.DB.NowFunc(),
	}

	if err := app.DB.Create(&log).Error; err != nil {
		return err
	}

	return nil
}

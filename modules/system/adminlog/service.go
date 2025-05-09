package adminlog

import (
	"fmt"

	"github.com/andycai/goapi/models"
	"github.com/gofiber/fiber/v2"
)

// WriteLog 创建操作日志
func WriteLog(c *fiber.Ctx, action string, resource string, resourceID uint, details string) error {
	currentUser := app.CurrentUser(c)

	if currentUser.ID == 0 {
		return fmt.Errorf("登录已过期，请重新登录")
	}

	log := models.AdminLog{
		UserID:     currentUser.ID,
		Username:   currentUser.Username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IP:         c.IP(),
		UserAgent:  c.Get("User-Agent"),
		CreatedAt:  app.DB.NowFunc(),
	}

	if err := app.DB.Create(&log).Error; err != nil {
		return err
	}

	return nil
}

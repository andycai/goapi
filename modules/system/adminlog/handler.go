package adminlog

import (
	"github.com/andycai/goapi/models"
	"github.com/gofiber/fiber/v2"
)

// listLogsHandler 获取操作日志列表
func listLogsHandler(c *fiber.Ctx) error {
	var total int64

	// 获取查询参数
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	username := c.Query("username")
	action := c.Query("action")
	resource := c.Query("resource")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// 构建查询
	query := app.DB.Model(&models.AdminLog{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	// 获取总数
	query.Count(&total)

	// 获取分页数据
	var logs []models.AdminLog
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "获取操作日志失败",
		})
	}

	return c.JSON(fiber.Map{
		"total": total,
		"data":  logs,
	})
}

// deleteLogsHandler 删除指定日期之前的操作日志
func deleteLogsHandler(c *fiber.Ctx) error {
	beforeDate := c.Query("beforeDate")
	if beforeDate == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "请指定日期",
		})
	}

	if err := app.DB.Where("created_at < ?", beforeDate).Delete(&models.AdminLog{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "删除操作日志失败",
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除成功",
	})
}

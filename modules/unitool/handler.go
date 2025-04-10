package unitool

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// findGuidHandler 查找重复GUID的请求处理器
func findGuidHandler(c *fiber.Ctx) error {
	var req FindGuidRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证请求数据
	if req.TargetPath == "" || req.NotificationURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "目标路径和通知URL不能为空",
		})
	}

	// 调用服务层处理
	log, err := FindDuplicateGuids(req.TargetPath, req.NotificationURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "查找GUID失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "已开始查找GUID",
		"log_id":  log.ID,
	})
}

// getFindGuidLogsHandler 获取GUID查找日志列表的请求处理器
func getFindGuidLogsHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	logs, total, err := GetFindGuidLogs(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取日志列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"logs":  logs,
		"total": total,
	})
}

// getDuplicateGuidsHandler 获取重复GUID列表的请求处理器
func getDuplicateGuidsHandler(c *fiber.Ctx) error {
	// 获取日志ID参数
	logID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的日志ID",
		})
	}

	duplicates, err := GetDuplicateGuids(uint(logID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取重复GUID列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"duplicates": duplicates,
	})
}

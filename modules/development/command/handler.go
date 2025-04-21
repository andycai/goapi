package command

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// getCommandsHandler 获取命令列表的请求处理器
func getCommandsHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	commands, total, err := GetCommands(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取命令列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"commands": commands,
		"total":    total,
	})
}

// createCommandHandler 创建命令的请求处理器
func createCommandHandler(c *fiber.Ctx) error {
	var cmd Command
	if err := c.BodyParser(&cmd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := CreateCommand(&cmd); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "创建命令失败: " + err.Error(),
		})
	}

	return c.JSON(cmd)
}

// executeCommandHandler 执行命令的请求处理器
func executeCommandHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的命令ID",
		})
	}

	execution, err := ExecuteCommand(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "执行命令失败: " + err.Error(),
		})
	}

	return c.JSON(execution)
}

// getCommandExecutionsHandler 获取命令执行记录的请求处理器
func getCommandExecutionsHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的命令ID",
		})
	}

	executions, err := GetCommandExecutions(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取执行记录失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"executions": executions,
	})
}

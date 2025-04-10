package parameter

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// listParametersHandler 获取参数列表
func listParametersHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	// 获取搜索参数
	search := c.Query("search", "")

	// 获取参数列表
	parameters, total, err := getParameters(limit, page, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取参数列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"parameters": parameters,
		"total":      total,
	})
}

// getParameterHandler 获取单个参数
func getParameterHandler(c *fiber.Ctx) error {
	// 获取参数ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的参数ID",
		})
	}

	// 获取参数
	parameter, err := getParameter(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取参数失败: " + err.Error(),
		})
	}

	return c.JSON(parameter)
}

// createParameterHandler 创建参数
func createParameterHandler(c *fiber.Ctx) error {
	// 解析请求
	req := new(ParameterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	currentUser := app.CurrentUser(c)

	// 创建参数
	parameter, err := createParameter(req, uint(currentUser.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "创建参数失败: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(parameter)
}

// updateParameterHandler 更新参数
func updateParameterHandler(c *fiber.Ctx) error {
	// 获取参数ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的参数ID",
		})
	}

	// 解析请求
	req := new(ParameterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	currentUser := app.CurrentUser(c)

	// 更新参数
	parameter, err := updateParameter(uint(id), req, uint(currentUser.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新参数失败: " + err.Error(),
		})
	}

	return c.JSON(parameter)
}

// deleteParameterHandler 删除参数
func deleteParameterHandler(c *fiber.Ctx) error {
	// 获取参数ID
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的参数ID",
		})
	}

	// 删除参数
	if err := deleteParameter(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除参数失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "参数删除成功",
	})
}

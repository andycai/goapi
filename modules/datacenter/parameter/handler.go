package parameter

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// listParameterHandler 获取参数列表
func listParameterHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	search := c.Query("search", "")

	parameters, total, err := QueryParameters(limit, page, search)
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
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的参数ID",
		})
	}

	parameter, err := QueryParameter(uint(id))
	if err != nil {
		if err == ErrParameterNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取参数失败: " + err.Error(),
		})
	}

	return c.JSON(parameter)
}

// addParameterHandler 添加参数
func addParameterHandler(c *fiber.Ctx) error {
	req := new(ParameterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	userID := c.Locals("userID").(uint)

	parameter, err := CommandCreateParameter(req, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "添加参数失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":   "添加参数成功",
		"parameter": parameter,
	})
}

// editParameterHandler 编辑参数
func editParameterHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的参数ID",
		})
	}

	req := new(ParameterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	userID := c.Locals("userID").(uint)

	parameter, err := CommandUpdateParameter(uint(id), req, userID)
	if err != nil {
		if err == ErrParameterNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新参数失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":   "更新参数成功",
		"parameter": parameter,
	})
}

// deleteParameterHandler 删除参数
func deleteParameterHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的参数ID",
		})
	}

	if err := CommandDeleteParameter(uint(id)); err != nil {
		if err == ErrParameterNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除参数失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除参数成功",
	})
}

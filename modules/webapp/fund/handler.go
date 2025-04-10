package fund

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// saveConfigHandler 保存配置
func saveConfigHandler(c *fiber.Ctx) error {
	config := new(FundConfig)
	if err := c.BodyParser(config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := updateConfig(config); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "保存配置失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "配置保存成功",
	})
}

// getConfigHandler 获取配置
func getConfigHandler(c *fiber.Ctx) error {
	return c.JSON(getConfig())
}

// syncDataHandler 同步数据
func syncDataHandler(c *fiber.Ctx) error {
	if err := SyncData(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "同步数据失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "数据同步成功",
	})
}

// listFundsHandler 获取基金列表
func listFundsHandler(c *fiber.Ctx) error {
	// 获取分页参数
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// 获取基金列表
	funds, total, err := ListFunds(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取基金列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"funds": funds,
		"total": total,
	})
}

// getMarketDataHandler 获取市场数据
func getMarketDataHandler(c *fiber.Ctx) error {
	indices, err := GetMarketIndices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取市场数据失败: " + err.Error(),
		})
	}

	if len(indices) == 0 {
		// 如果没有数据，返回空数组而不是null
		return c.JSON(fiber.Map{
			"indices": []MarketIndex{},
		})
	}

	return c.JSON(fiber.Map{
		"indices": indices,
	})
}

// getHotFundsHandler 获取热门基金
func getHotFundsHandler(c *fiber.Ctx) error {
	// 获取前10只热门基金
	funds, _, err := ListFunds(1, 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取热门基金失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"funds": funds,
	})
}

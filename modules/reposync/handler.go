package reposync

import (
	"strconv"

	"github.com/andycai/unitool/core/utility/path"
	"github.com/gofiber/fiber/v2"
)

// saveConfigHandler 保存仓库配置
func saveConfigHandler(c *fiber.Ctx) error {
	config := new(RepoConfig)
	if err := c.BodyParser(config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证路径安全性
	if !path.IsValid(config.LocalPath1) || !path.IsValid(config.LocalPath2) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的本地路径",
		})
	}

	if err := updateConfig(config); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "保存配置失败",
		})
	}

	return c.JSON(fiber.Map{
		"message": "配置保存成功",
	})
}

// getConfigHandler 获取仓库配置
func getConfigHandler(c *fiber.Ctx) error {
	config := getConfig()
	return c.JSON(config)
}

// checkoutHandler 检出仓库
func checkoutHandler(c *fiber.Ctx) error {
	if err := checkoutRepos(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "检出仓库失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "仓库检出成功",
	})
}

// listCommitsHandler 获取提交记录列表
func listCommitsHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	commits, totalCount, err := getCommits(limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取提交记录失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"commits": commits,
		"total":   totalCount,
	})
}

// syncCommitsHandler 同步提交记录
func syncCommitsHandler(c *fiber.Ctx) error {
	type SyncRequest struct {
		Revisions []string `json:"revisions"`
	}

	var req SyncRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if len(req.Revisions) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "版本号列表不能为空",
		})
	}

	if err := SyncCommits(req.Revisions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "同步提交失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "同步成功",
	})
}

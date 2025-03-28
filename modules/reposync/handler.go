package reposync

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ConfigHandler 保存仓库配置
func ConfigHandler(c *fiber.Ctx) error {
	config := new(RepoConfig)
	if err := c.BodyParser(config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证路径安全性
	if !GetService().isValidPath(config.LocalPath1) || !GetService().isValidPath(config.LocalPath2) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的本地路径",
		})
	}

	if err := GetService().UpdateConfig(config); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "保存配置失败",
		})
	}

	return c.JSON(fiber.Map{
		"message": "配置保存成功",
	})
}

// GetConfigHandler 获取仓库配置
func GetConfigHandler(c *fiber.Ctx) error {
	config := GetService().GetConfig()
	return c.JSON(config)
}

// CheckoutHandler 检出仓库
func CheckoutHandler(c *fiber.Ctx) error {
	if err := GetService().CheckoutRepos(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "检出仓库失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "仓库检出成功",
	})
}

// ListCommitsHandler 获取提交记录列表
func ListCommitsHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	commits, totalCount, err := GetService().GetCommits(limit, page)
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

// SyncHandler 同步提交记录
func SyncHandler(c *fiber.Ctx) error {
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

	if err := GetService().SyncCommits(req.Revisions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "同步提交失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "同步成功",
	})
}

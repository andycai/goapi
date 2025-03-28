package reposync

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// saveConfig 保存仓库配置
func saveConfig(c *fiber.Ctx) error {
	config := new(RepoConfig)
	if err := c.BodyParser(config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证路径安全性
	if !srv.isValidPath(config.LocalPath1) || !srv.isValidPath(config.LocalPath2) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的本地路径",
		})
	}

	if err := srv.UpdateConfig(config); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "保存配置失败",
		})
	}

	return c.JSON(fiber.Map{
		"message": "配置保存成功",
	})
}

// getConfig 获取仓库配置
func getConfig(c *fiber.Ctx) error {
	config := srv.GetConfig()
	return c.JSON(config)
}

// checkout 检出仓库
func checkout(c *fiber.Ctx) error {
	if err := srv.CheckoutRepos(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "检出仓库失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "仓库检出成功",
	})
}

// getCommits 获取提交记录列表
func getCommits(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	commits, totalCount, err := srv.GetCommits(limit, page)
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

// syncCommits 同步提交记录
func syncCommits(c *fiber.Ctx) error {
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

	if err := srv.SyncCommits(req.Revisions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "同步提交失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "同步成功",
	})
}

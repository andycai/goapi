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

// syncPublicCommitsHandler 公开的同步API，同步两个版本间的差异
func syncPublicCommitsHandler(c *fiber.Ctx) error {
	var req struct {
		FromRevision string `json:"fromRevision"`
		ToRevision   string `json:"toRevision"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if req.FromRevision == "" || req.ToRevision == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fromRevision 和 toRevision 不能为空",
		})
	}

	// 检查版本号大小关系
	// 尝试将版本号转换为数字进行比较
	fromNum, fromErr := strconv.Atoi(req.FromRevision)
	toNum, toErr := strconv.Atoi(req.ToRevision)

	if fromErr == nil && toErr == nil {
		// 如果都是数字，直接比较数值
		if fromNum >= toNum {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "fromRevision 必须小于 toRevision",
			})
		}
	} else {
		// 如果不是数字，尝试按字符串比较
		// 对于某些版本号形式，字符串比较可能不准确
		if req.FromRevision >= req.ToRevision {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "fromRevision 必须小于 toRevision",
			})
		}
	}

	changeCount, err := SyncChangesBetweenRevisions(req.FromRevision, req.ToRevision)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "同步成功",
		"changes": changeCount,
	})
}

// syncPublicAutoHandler 自动同步未同步的提交记录
func syncPublicAutoHandler(c *fiber.Ctx) error {
	fromRev, toRev, err := FindUnsyncedRevisionRange()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 如果没有找到未同步的提交，返回成功但无变更
	if fromRev == "" || toRev == "" {
		return c.JSON(fiber.Map{
			"message": "没有需要同步的提交",
			"changes": 0,
		})
	}

	changeCount, err := SyncChangesBetweenRevisions(fromRev, toRev)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "同步成功",
		"changes": changeCount,
		"from":    fromRev,
		"to":      toRev,
	})
}

// refreshCommitsHandler 刷新提交记录
func refreshCommitsHandler(c *fiber.Ctx) error {
	var req struct {
		Limit int `json:"limit"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if req.Limit <= 0 || req.Limit > 1000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "limit 必须在 1-1000 之间",
		})
	}

	if err := RefreshCommits(req.Limit); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "刷新提交记录失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "刷新提交记录成功",
	})
}

// clearSyncDataHandler 清空同步数据
func clearSyncDataHandler(c *fiber.Ctx) error {
	if err := ClearSyncData(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "清空数据失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "清空数据成功",
	})
}
